package dbseeder

import (
	"context"
	"log/slog"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Seeder struct {
	queries *repo.Queries
}

func NewSeeder(queries *repo.Queries) *Seeder {
	return &Seeder{queries: queries}
}

// helper function to extract ID from a row, with fallback to array index
func extractID(row map[string]interface{}, rowIndex int) int32 {
	if id, ok := row["id"]; ok {
		// Handle multiple types that JSON might provide
		switch typedVal := id.(type) {
		case float64:
			return int32(typedVal)
		case int32:
			return typedVal
		case int64:
			return int32(typedVal)
		case int:
			return int32(typedVal)
		}
	}
	return int32(rowIndex)
}

// helper function to extract int32 value from any numeric type
func extractInt32(value interface{}, defaultValue int32) int32 {
	switch typedVal := value.(type) {
	case float64:
		return int32(typedVal)
	case int32:
		return typedVal
	case int64:
		return int32(typedVal)
	case int:
		return int32(typedVal)
	case pgtype.Int4:
		if typedVal.Valid {
			return typedVal.Int32
		}
		return defaultValue
	default:
		return defaultValue
	}
}

// IsAlreadySeeded checks if the database has been seeded with data
func (s *Seeder) IsAlreadySeeded(ctx context.Context) (bool, error) {
	statsCount, err := s.queries.GetStatsCount(ctx)
	if err != nil {
		return false, err
	}
	return statsCount > 0, nil
}

func (s *Seeder) isDatabaseSeeded(ctx context.Context) (bool, error) {
	return s.IsAlreadySeeded(ctx)
}
func (s *Seeder) SeedDatabase(ctx context.Context, data *TransformedData) error {
	slog.Info("Starting database seeding")

	// Check if database is already seeded
	seeded, err := s.isDatabaseSeeded(ctx)
	if err != nil {
		// Log the error but continue (database might be empty)
		slog.Debug("Error checking if database is seeded", "error", err)
	} else if seeded {
		slog.Info("Database already seeded, skipping seeding process")
		return nil
	}

	// Order matters: insert tables with no foreign keys first
	if err := s.seedStats(ctx, data); err != nil {
		return err
	}

	if err := s.seedClasses(ctx, data); err != nil {
		return err
	}

	if err := s.seedRaces(ctx, data); err != nil {
		return err
	}

	if err := s.seedMaterials(ctx, data); err != nil {
		return err
	}

	if err := s.seedMaterialStats(ctx, data); err != nil {
		return err
	}

	if err := s.seedTraits(ctx, data); err != nil {
		return err
	}

	if err := s.seedCreatures(ctx, data); err != nil {
		return err
	}

	if err := s.seedSpecializations(ctx, data); err != nil {
		return err
	}

	if err := s.seedPerks(ctx, data); err != nil {
		return err
	}

	if err := s.seedArtifacts(ctx, data); err != nil {
		return err
	}

	if err := s.seedSpells(ctx, data); err != nil {
		return err
	}

	if err := s.seedSpellProperties(ctx, data); err != nil {
		return err
	}

	slog.Info("Database seeding completed successfully")
	return nil
}

func (s *Seeder) seedStats(ctx context.Context, data *TransformedData) error {
	statsTable, exists := data.GetTable("stats")
	if !exists {
		slog.Debug("Stats table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range statsTable {
		statType, ok := row["stat_type"].(string)
		if !ok {
			slog.Warn("Invalid stat type in row", "row", row)
			continue
		}

		id := extractID(row, rowID)
		_, err := s.queries.CreateStat(ctx, repo.CreateStatParams{
			ID:   id,
			Type: repo.StatType(statType),
		})
		if err != nil {
			slog.Error("Failed to create stat", "statType", statType, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded stats", "count", successCount)
	return nil
}

func (s *Seeder) seedClasses(ctx context.Context, data *TransformedData) error {
	classesTable, exists := data.GetTable("classes")
	if !exists {
		slog.Debug("Classes table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range classesTable {
		name, ok := row["name"].(string)
		if !ok {
			slog.Warn("Invalid class name in row", "row", row)
			continue
		}

		icon, ok := row["icon"].([]byte)
		if !ok {
			icon = nil
		}

		id := extractID(row, rowID)
		_, err := s.queries.CreateClass(ctx, repo.CreateClassParams{
			ID:   id,
			Name: name,
			Icon: icon,
		})
		if err != nil {
			slog.Error("Failed to create class", "name", name, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded classes", "count", successCount)
	return nil
}

func (s *Seeder) seedRaces(ctx context.Context, data *TransformedData) error {
	racesTable, exists := data.GetTable("races")
	if !exists {
		slog.Debug("Races table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range racesTable {
		name, ok := row["name"].(string)
		if !ok {
			slog.Warn("Invalid race name in row", "row", row)
			continue
		}

		icon, ok := row["icon"].([]byte)
		if !ok {
			icon = nil
		}

		id := extractID(row, rowID)
		_, err := s.queries.CreateRace(ctx, repo.CreateRaceParams{
			ID:   id,
			Name: name,
			Icon: icon,
		})
		if err != nil {
			slog.Error("Failed to create race", "name", name, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded races", "count", successCount)
	return nil
}

func (s *Seeder) seedMaterials(ctx context.Context, data *TransformedData) error {
	materialsTable, exists := data.GetTable("materials")
	if !exists {
		slog.Debug("Materials table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range materialsTable {
		name, ok := row["name"].(string)
		if !ok {
			slog.Warn("Invalid material name in row", "row", row)
			continue
		}

		icon, ok := row["icon"].([]byte)
		if !ok {
			icon = []byte{}
		}

		materialType, ok := row["type"].(string)
		if !ok {
			materialType = ""
		}

		id := extractID(row, rowID)
		_, err := s.queries.CreateMaterial(ctx, repo.CreateMaterialParams{
			ID:   id,
			Name: name,
			Icon: icon,
			Type: repo.NullMaterialType{
				MaterialType: repo.MaterialType(materialType),
				Valid:        materialType != "",
			},
		})
		if err != nil {
			slog.Error("Failed to create material", "name", name, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded materials", "count", successCount)
	return nil
}

func (s *Seeder) seedMaterialStats(ctx context.Context, data *TransformedData) error {
	materialStatsTable, exists := data.GetTable("material_stats")
	if !exists {
		slog.Debug("Material stats table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range materialStatsTable {
		materialID := int32(0)
		if mID, ok := row["material_id"]; ok {
			materialID = extractInt32(mID, 0)
		}

		if materialID < 0 {
			slog.Warn("Invalid material_id in row", "row", row)
			continue
		}

		statID := int32(0)
		if sID, ok := row["stat_id"]; ok {
			statID = extractInt32(sID, 0)
		}

		if statID < 0 {
			slog.Warn("Invalid stat_id in row", "row", row)
			continue
		}

		statId2 := pgtype.Int4{Valid: false}
		if sID2, ok := row["stat_id2"]; ok && sID2 != nil {
			id := extractInt32(sID2, 0)
			if id > 0 {
				statId2 = pgtype.Int4{Int32: id, Valid: true}
			}
		}

		id := extractID(row, rowID)
		_, err := s.queries.CreateMaterialStat(ctx, repo.CreateMaterialStatParams{
			ID:         id,
			MaterialID: materialID,
			StatID:     statID,
			StatId2:    statId2,
		})
		if err != nil {
			slog.Error("Failed to create material stat", "materialID", materialID, "statID", statID, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded material stats", "count", successCount)
	return nil
}

func (s *Seeder) seedTraits(ctx context.Context, data *TransformedData) error {
	traitsTable, exists := data.GetTable("traits")
	if !exists {
		slog.Debug("Traits table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range traitsTable {
		name, ok := row["name"].(string)
		if !ok {
			slog.Warn("Invalid trait name in row", "row", row)
			continue
		}

		description, _ := row["description"].(string)

		materialID := pgtype.Int4{Valid: false}
		if matID, ok := row["material_id"]; ok && matID != nil {
			id := extractInt32(matID, 0)
			if id > 0 {
				// Verify material exists in database before setting foreign key
				exists, err := s.queries.MaterialExists(ctx, id)
				if err != nil {
					slog.Warn("Failed to check if material exists", "materialID", id, "error", err)
				} else if exists {
					materialID = pgtype.Int4{Int32: id, Valid: true}
				} else {
					slog.Debug("Material not found, skipping material_id", "traitName", name, "materialID", id)
				}
			}
		}

		id := extractID(row, rowID)
		_, err := s.queries.CreateTrait(ctx, repo.CreateTraitParams{
			ID:          id,
			Name:        name,
			Description: description,
			MaterialID:  materialID,
		})
		if err != nil {
			slog.Error("Failed to create trait", "name", name, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded traits", "count", successCount)
	return nil
}

func (s *Seeder) seedCreatures(ctx context.Context, data *TransformedData) error {
	creaturesTable, exists := data.GetTable("creatures")
	if !exists {
		slog.Debug("Creatures table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range creaturesTable {
		name, ok := row["name"].(string)
		if !ok {
			slog.Warn("Invalid creature name in row", "row", row)
			continue
		}

		image, ok := row["image"].([]byte)
		if !ok {
			image = nil
		}

		traitID := pgtype.Int4{Valid: false}
		if tID, ok := row["trait_id"]; ok {
			id := extractInt32(tID, 0)
			if id > 0 {
				// Verify trait exists in database before setting foreign key
				exists, err := s.queries.TraitExists(ctx, id)
				if err != nil {
					slog.Warn("Failed to check if trait exists", "traitID", id, "error", err)
				} else if exists {
					traitID = pgtype.Int4{Int32: id, Valid: true}
				} else {
					slog.Debug("Trait not found, skipping trait_id", "creatureName", name, "traitID", id)
				}
			}
		}

		classID := int32(0)
		if cID, ok := row["class_id"]; ok {
			classID = extractInt32(cID, 0)
		}

		raceID := int32(0)
		if rID, ok := row["race_id"]; ok {
			raceID = extractInt32(rID, 0)
		}

		id := extractID(row, rowID)
		_, err := s.queries.CreateCreature(ctx, repo.CreateCreatureParams{
			ID:      id,
			Name:    name,
			Image:   image,
			TraitID: traitID,
			ClassID: classID,
			RaceID:  raceID,
		})
		if err != nil {
			slog.Error("Failed to create creature", "name", name, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded creatures", "count", successCount)
	return nil
}

func (s *Seeder) seedSpecializations(ctx context.Context, data *TransformedData) error {
	specializationsTable, exists := data.GetTable("specializations")
	if !exists {
		slog.Debug("Specializations table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range specializationsTable {
		name, ok := row["name"].(string)
		if !ok {
			slog.Warn("Invalid specialization name in row", "row", row)
			continue
		}

		description, _ := row["description"].(string)

		id := extractID(row, rowID)
		_, err := s.queries.CreateSpecialization(ctx, repo.CreateSpecializationParams{
			ID:          id,
			Name:        name,
			Description: description,
		})
		if err != nil {
			slog.Error("Failed to create specialization", "name", name, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded specializations", "count", successCount)
	return nil
}

func (s *Seeder) seedPerks(ctx context.Context, data *TransformedData) error {
	perksTable, exists := data.GetTable("perks")
	if !exists {
		slog.Debug("Perks table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range perksTable {
		name, ok := row["name"].(string)
		if !ok {
			slog.Warn("Invalid perk name in row", "row", row)
			continue
		}

		description, _ := row["description"].(string)
		icon, _ := row["icon"].([]byte)

		specializationID := int32(0)
		if sID, ok := row["specialization_id"]; ok {
			specializationID = extractInt32(sID, 0)
		}

		id := extractID(row, rowID)
		_, err := s.queries.CreatePerk(ctx, repo.CreatePerkParams{
			ID:               id,
			Name:             name,
			Description:      description,
			Icon:             icon,
			SpecializationID: specializationID,
		})
		if err != nil {
			slog.Error("Failed to create perk", "name", name, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded perks", "count", successCount)
	return nil
}

func (s *Seeder) seedArtifacts(ctx context.Context, data *TransformedData) error {
	artifactsTable, exists := data.GetTable("artifacts")
	if !exists {
		slog.Debug("Artifacts table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range artifactsTable {
		name, ok := row["name"].(string)
		if !ok {
			slog.Warn("Invalid artifact name in row", "row", row)
			continue
		}

		icon, _ := row["icon"].([]byte)

		artifactType, ok := row["stat_type"].(string)
		if !ok {
			artifactType = ""
		}

		id := extractID(row, rowID)
		_, err := s.queries.CreateArtifact(ctx, repo.CreateArtifactParams{
			ID:   id,
			Name: name,
			Icon: icon,
			Type: repo.StatType(artifactType),
		})
		if err != nil {
			slog.Error("Failed to create artifact", "name", name, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded artifacts", "count", successCount)
	return nil
}

func (s *Seeder) seedSpells(ctx context.Context, data *TransformedData) error {
	spellsTable, exists := data.GetTable("spells")
	if !exists {
		slog.Debug("Spells table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range spellsTable {
		name, ok := row["name"].(string)
		if !ok {
			slog.Warn("Invalid spell name in row", "row", row)
			continue
		}

		description, _ := row["description"].(string)
		icon, _ := row["icon"].([]byte)

		charges := int32(0)
		if c, ok := row["charges"]; ok {
			charges = extractInt32(c, 0)
		}

		classID := int32(0)
		if cID, ok := row["class_id"]; ok {
			classID = extractInt32(cID, 0)
		}

		id := extractID(row, rowID)
		_, err := s.queries.CreateSpell(ctx, repo.CreateSpellParams{
			ID:          id,
			Name:        name,
			Description: description,
			Icon:        icon,
			Charges:     charges,
			ClassID:     classID,
		})
		if err != nil {
			slog.Error("Failed to create spell", "name", name, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded spells", "count", successCount)
	return nil
}

func (s *Seeder) seedSpellProperties(ctx context.Context, data *TransformedData) error {
	spellPropertiesTable, exists := data.GetTable("spellProperties")
	if !exists {
		slog.Debug("Spell properties table not found in transformed data")
		return nil
	}

	successCount := 0
	for rowID, row := range spellPropertiesTable {
		name, ok := row["name"].(string)
		if !ok {
			slog.Warn("Invalid spell property name in row", "row", row)
			continue
		}

		materialID := int32(0)
		if mID, ok := row["material_id"]; ok {
			materialID = extractInt32(mID, 0)
		}

		id := extractID(row, rowID)
		_, err := s.queries.CreateSpellProperty(ctx, repo.CreateSpellPropertyParams{
			ID:         id,
			Name:       name,
			MaterialID: materialID,
		})
		if err != nil {
			slog.Error("Failed to create spell property", "name", name, "error", err)
			return err
		}
		successCount++
	}

	slog.Info("Seeded spell properties", "count", successCount)
	return nil
}
