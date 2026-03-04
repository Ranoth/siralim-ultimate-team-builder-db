package dbseeder

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type inserter struct {
	logger  *slog.Logger
	config  *config
	queries *repo.Queries
}

func newInserter(logger *slog.Logger, config *config, queries *repo.Queries) *inserter {
	return &inserter{logger: logger, config: config, queries: queries}
}

func batchInsertTable[T any](
	i *inserter,
	ctx context.Context,
	tableName string,
	insertFunc func(context.Context, []T) (int64, error),
) error {
	items := i.config.jsonSources[tableName].items
	if len(items) == 0 {
		return nil
	}

	params := make([]T, 0, len(items))
	for _, item := range items {
		param, err := mapToTyped[T](item)
		if err != nil {
			i.logger.Error("Failed to convert map to typed struct", "table", tableName, "error", err)
			continue
		}
		params = append(params, *param)
	}

	if len(params) == 0 {
		return nil
	}

	count, err := insertFunc(ctx, params)
	if err != nil {
		i.logger.Error("Failed to batch insert", "table", tableName, "error", err)
		return err
	}

	i.logger.Info("Batch inserted records", "table", tableName, "count", count)
	return nil
}

func mapToTyped[T any](data map[string]interface{}) (*T, error) {
	normalizedData := make(map[string]interface{}, len(data))
	for key, value := range data {
		normalizedData[key] = value
	}

	if err := normalizeNullableInt4Fields[T](normalizedData); err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(normalizedData)
	if err != nil {
		return nil, err
	}
	var result T
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}

	if err := applyNullableInt4Fields(&result, normalizedData); err != nil {
		return nil, err
	}

	return &result, nil
}

var pgtypeInt4Type = reflect.TypeOf(pgtype.Int4{})

func normalizeNullableInt4Fields[T any](data map[string]interface{}) error {
	typedValue := reflect.TypeOf((*T)(nil)).Elem()

	for i := range typedValue.NumField() {
		field := typedValue.Field(i)
		if field.Type != pgtypeInt4Type {
			continue
		}

		jsonFieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			jsonFieldName = strings.Split(jsonTag, ",")[0]
		}

		rawValue, exists := data[jsonFieldName]
		if !exists {
			continue
		}

		nullableInt4, err := toNullableInt4(rawValue)
		if err != nil {
			return fmt.Errorf("failed to parse nullable int4 field %q: %w", jsonFieldName, err)
		}

		data[jsonFieldName] = nullableInt4
	}

	return nil
}

func toNullableInt4(rawValue interface{}) (pgtype.Int4, error) {
	if rawValue == nil {
		return pgtype.Int4{Valid: false}, nil
	}

	switch value := rawValue.(type) {
	case pgtype.Int4:
		return value, nil
	case int:
		return pgtype.Int4{Int32: int32(value), Valid: true}, nil
	case int8:
		return pgtype.Int4{Int32: int32(value), Valid: true}, nil
	case int16:
		return pgtype.Int4{Int32: int32(value), Valid: true}, nil
	case int32:
		return pgtype.Int4{Int32: value, Valid: true}, nil
	case int64:
		return pgtype.Int4{Int32: int32(value), Valid: true}, nil
	case float32:
		if value != float32(int32(value)) {
			return pgtype.Int4{}, fmt.Errorf("non-integer float value %v", value)
		}
		return pgtype.Int4{Int32: int32(value), Valid: true}, nil
	case float64:
		if value != float64(int32(value)) {
			return pgtype.Int4{}, fmt.Errorf("non-integer float value %v", value)
		}
		return pgtype.Int4{Int32: int32(value), Valid: true}, nil
	case string:
		if value == "" {
			return pgtype.Int4{Valid: false}, nil
		}
		parsed, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return pgtype.Int4{}, fmt.Errorf("invalid integer string %q", value)
		}
		return pgtype.Int4{Int32: int32(parsed), Valid: true}, nil
	default:
		return pgtype.Int4{}, fmt.Errorf("unsupported type %T", rawValue)
	}
}

func applyNullableInt4Fields[T any](target *T, data map[string]interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}

	structValue := targetValue.Elem()
	structType := structValue.Type()

	for i := range structType.NumField() {
		fieldType := structType.Field(i)
		if fieldType.Type != pgtypeInt4Type {
			continue
		}

		jsonFieldName := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			jsonFieldName = strings.Split(jsonTag, ",")[0]
		}

		rawValue, exists := data[jsonFieldName]
		if !exists {
			structValue.Field(i).Set(reflect.ValueOf(pgtype.Int4{Valid: false}))
			continue
		}

		nullableValue, err := toNullableInt4(rawValue)
		if err != nil {
			return fmt.Errorf("failed to parse nullable int4 field %q: %w", jsonFieldName, err)
		}

		structValue.Field(i).Set(reflect.ValueOf(nullableValue))
	}

	return nil
}

func (i *inserter) insert() {
	ctx := context.Background()

	// Insert in dependency order
	insertOrder := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"classes", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "classes", i.queries.BatchInsertClasses)
		}},
		{"stats", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "stats", i.queries.BatchInsertStats)
		}},
		{"races", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "races", i.queries.BatchInsertRaces)
		}},
		{"specializations", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "specializations", i.queries.BatchInsertSpecializations)
		}},
		{"materials", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "materials", i.queries.BatchInsertMaterials)
		}},
		{"material_stats", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "material_stats", i.queries.BatchInsertMaterialStats)
		}},
		{"traits", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "traits", i.queries.BatchInsertTraits)
		}},
		{"artifacts", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "artifacts", i.queries.BatchInsertArtifacts)
		}},
		{"spellProperties", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "spellProperties", i.queries.BatchInsertSpellProperties)
		}},
		{"perks", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "perks", i.queries.BatchInsertPerks)
		}},
		{"spells", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "spells", i.queries.BatchInsertSpells)
		}},
		{"creatures", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "creatures", i.queries.BatchInsertCreatures)
		}},
		{"relics", func(ctx context.Context) error {
			return batchInsertTable(i, ctx, "relics", i.queries.BatchInsertRelics)
		}},
	}

	for _, insertOp := range insertOrder {
		if err := insertOp.fn(ctx); err != nil {
			i.logger.Error("Failed to insert table", "table", insertOp.name, "error", err)
		}
	}
}
