package dbseeder

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

func FindIndexByValue(array []string, value string) int {
	for i, v := range array {
		if v == value {
			return i
		}
	}
	return -1
}

func GetIDValue(value interface{}) (int, bool) {
	switch typedValue := value.(type) {
	case float64:
		return int(typedValue), true
	case int:
		return typedValue, true
	case int32:
		return int(typedValue), true
	case int64:
		return int(typedValue), true
	case string:
		parsed, err := strconv.Atoi(typedValue)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func GetFieldValue(jsonMeta jsonMeta, jsonFieldName string, rowIndex int) (interface{}, bool) {
	if jsonFieldName == "" {
		return nil, false
	}

	fieldData, exists := jsonMeta.Fields[jsonFieldName]
	if !exists {
		return nil, false
	}

	value, exists := fieldData[rowIndex]
	return value, exists
}

func ArtifactIconPathFromValue(rawValue interface{}) (string, bool) {
	const artifactIconIndex = 5
	switch iconPaths := rawValue.(type) {
	case []interface{}:
		if len(iconPaths) <= artifactIconIndex {
			return "", false
		}
		iconPath, ok := iconPaths[artifactIconIndex].(string)
		if !ok || iconPath == "" {
			return "", false
		}
		return iconPath, true
	case []string:
		if len(iconPaths) <= artifactIconIndex || iconPaths[artifactIconIndex] == "" {
			return "", false
		}
		return iconPaths[artifactIconIndex], true
	default:
		return "", false
	}
}

func IconToBytes(iconPath string) ([]byte, bool) {
	absolutePath := filepath.Join(gameDataRootPath, iconPath)
	iconBytes, err := os.ReadFile(absolutePath)
	if err != nil {
		slog.Warn("Skipping icon conversion", "iconPath", absolutePath, "error", err)
		return nil, false
	}

	return iconBytes, true
}

func TransformIconsForDB(tableName, dbFieldName string, rawValue interface{}) interface{} {
	if tableName == "artifacts" && dbFieldName == "icon" {
		if iconPath, ok := ArtifactIconPathFromValue(rawValue); ok {
			if iconBytes, ok := IconToBytes(iconPath); ok {
				return iconBytes
			}
		}
		slog.Warn("Artifact icon conversion failed; storing empty bytes")
		return []byte{}
	}
	if tableName != "artifacts" && dbFieldName == "icon" {
		if iconPath, ok := rawValue.(string); ok && iconPath != "" {
			if iconBytes, ok := IconToBytes(iconPath); ok {
				return iconBytes
			}
		}
		slog.Warn("Skipping icon conversion for table", "tableName", tableName, "fieldName", dbFieldName)
		return nil
	}
	return rawValue
}

func TransformClassForDB(tableName string, dbFieldName string, rawValue interface{}) interface{} {
	if dbFieldName == "class_id" {
		if className, ok := rawValue.(string); ok {
			classID, exists := classes[className]
			if !exists {
				slog.Warn("Class name not found in mapping", "className", className)
				return nil
			}
			return classID.ID
		}
		slog.Warn("Invalid class name format", "rawValue", rawValue)
		return nil
	}
	return rawValue
}

func TransformRaceForDB(tableName string, dbFieldName string, rawValue interface{}, result *TransformedData) interface{} {
	if dbFieldName == "race_id" {
		if raceName, ok := rawValue.(string); ok {
			racesTable, exists := result.GetTable("races")
			if !exists {
				slog.Warn("Races table not found")
				return nil
			}
			for id, fields := range racesTable {
				if name, ok := fields["name"].(string); ok && name == raceName {
					return id
				}
			}
			slog.Warn("Race name not found in races table", "raceName", raceName)
			return nil
		}
		slog.Warn("Invalid race name format", "rawValue", rawValue)
		return nil
	}
	return rawValue
}

func TransformValueForDB(tableName string, dbFieldName string, rawValue interface{}, result *TransformedData) interface{} {
	rawValue = TransformIconsForDB(tableName, dbFieldName, rawValue)
	rawValue = TransformClassForDB(tableName, dbFieldName, rawValue)
	rawValue = TransformRaceForDB(tableName, dbFieldName, rawValue, result)
	return rawValue
}

func CorrelateFieldsForRow(jsonMeta jsonMeta, rowIndex int, mappings []fieldMapping, result *TransformedData) TableRow {
	rowData := make(TableRow)

	for _, mapping := range mappings {
		jsonFieldName := mapping.JSONField
		dbFieldName := mapping.DBField

		if value, exists := GetFieldValue(jsonMeta, jsonFieldName, rowIndex); exists {
			rowData[dbFieldName] = TransformValueForDB(jsonMeta.Name, dbFieldName, value, result)
		}
	}

	return rowData
}

func ProcessTableRow(jsonMeta jsonMeta, rowIndex int, idValue interface{}, result *TransformedData) {
	id, ok := GetIDValue(idValue)
	if !ok {
		slog.Warn("Skipping row due to missing or invalid id", "name", jsonMeta.Name, "id", idValue)
		return
	}

	fieldMapping, exists := correlatedFieldNamesMetaMap[jsonMeta.Name]
	if !exists {
		slog.Warn("No field mapping found for table", "name", jsonMeta.Name)
		return
	}

	// Ensure table and row exist
	table := result.EnsureTable(jsonMeta.Name)
	if table[id] == nil {
		table[id] = make(TableRow)
	}

	rowData := CorrelateFieldsForRow(jsonMeta, rowIndex, fieldMapping, result)
	for field, value := range rowData {
		table[id][field] = value
	}
}

func ProcessTableData(jsonMeta jsonMeta, result *TransformedData) {
	result.EnsureTable(jsonMeta.Name)

	idField, exists := jsonMeta.Fields["id"]
	if !exists {
		slog.Warn("No id field found for table", "name", jsonMeta.Name)
		idField = make(map[int]interface{})

		for _, fieldData := range jsonMeta.Fields {
			counter := 0
			for rowIndex := range fieldData {
				idField[rowIndex] = counter
				counter++
			}
			break
		}
		slog.Info("Generated synthetic id field for table", "name", jsonMeta.Name, "rows", len(idField))
	}

	for rowIndex, idValue := range idField {
		ProcessTableRow(jsonMeta, rowIndex, idValue, result)
	}
}

func PopulateStaticTable(tableName string, sourceData map[string]staticTableEntry, result *TransformedData) {
	table := result.EnsureTable(tableName)

	for itemName, itemData := range sourceData {
		id := itemData.ID
		iconPath := itemData.IconPath

		row := make(TableRow)
		row["id"] = id

		// Use field name based on table type
		if tableName == "stats" {
			row["stat_type"] = itemName
		} else {
			row["name"] = itemName
		}

		if iconBytes, ok := IconToBytes(iconPath); ok {
			row["icon"] = iconBytes
		} else {
			slog.Warn("Failed to read icon for static table", "tableName", tableName, "itemName", itemName, "iconPath", iconPath)
			row["icon"] = nil
		}

		table[id] = row
	}

	slog.Info("Populated static table", "tableName", tableName, "count", len(table))
}
