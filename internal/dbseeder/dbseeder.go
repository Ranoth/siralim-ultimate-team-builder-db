package dbseeder

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

type jsonMeta struct {
	FilePath string
	Name     string
	Fields   map[string]map[int]interface{}
}

var classes = map[string][2]interface{}{
	"life":    {0, "images/misc/class/life.png"},
	"death":   {1, "images/misc/class/death.png"},
	"nature":  {2, "images/misc/class/nature.png"},
	"sorcery": {3, "images/misc/class/sorcery.png"},
	"chaos":   {4, "images/misc/class/chaos.png"},
}
var stats = map[string][2]interface{}{
	"health":       {0, "images/misc/stat/health.png"},
	"attack":       {1, "images/misc/stat/attack.png"},
	"intelligence": {2, "images/misc/stat/intelligence.png"},
	"defense":      {3, "images/misc/stat/defense.png"},
	"speed":        {4, "images/misc/stat/speed.png"},
}

var jsonMetaArray = [9]jsonMeta{
	{FilePath: "/app/gameData/artifacts.json", Name: "artifacts"},
	{FilePath: "/app/gameData/races.json", Name: "races"},
	{FilePath: "/app/gameData/specializations.json", Name: "specializations"},
	{FilePath: "/app/gameData/materials.json", Name: "materials"},
	{FilePath: "/app/gameData/traits.json", Name: "traits"},
	{FilePath: "/app/gameData/creatures.json", Name: "creatures"},
	{FilePath: "/app/gameData/perks.json", Name: "perks"},
	{FilePath: "/app/gameData/spells.json", Name: "spells"},
	{FilePath: "/app/gameData/spellProperties.json", Name: "spellProperties"},
}

var correllatedFieldNamesMetaMap = map[string]map[string]map[string]string{
	"artifacts": {
		"id":        {"db": "id", "json": "id"},
		"name":      {"db": "name", "json": "name"},
		"icon":      {"db": "icon", "json": "icons"},
		"stat_type": {"db": "stat_type", "json": "stat"},
	},
	"classes": {
		"id":   {"db": "id", "json": ""},
		"name": {"db": "name", "json": ""},
		"icon": {"db": "icon", "json": ""},
	},
	"races": {
		"id":   {"db": "id", "json": ""},
		"name": {"db": "name", "json": "name"},
		"icon": {"db": "icon", "json": "icon"},
	},
	"specializations": {
		"id":          {"db": "id", "json": "id"},
		"name":        {"db": "name", "json": "name"},
		"description": {"db": "description", "json": "description"},
		"icon":        {"db": "icon", "json": "icon"},
	},
	"materials": {
		"id":          {"db": "id", "json": "id"},
		"name":        {"db": "name", "json": "name"},
		"description": {"db": "description", "json": "description"},
		"icon":        {"db": "icon", "json": "icon"},
		"type":        {"db": "type", "json": "slot"},
	},
	"traits": {
		"id":          {"db": "id", "json": "id"},
		"name":        {"db": "name", "json": "name"},
		"description": {"db": "description", "json": "description"},
		"material_id": {"db": "material_id", "json": "materialId"},
	},
	"creatures": {
		"id":          {"db": "id", "json": "id"},
		"name":        {"db": "name", "json": "name"},
		"description": {"db": "description", "json": "description"},
		"icon":        {"db": "icon", "json": "icon"},
		"trait_id":    {"db": "trait_id", "json": "trait"},
		"class_id":    {"db": "class_id", "json": "class"},
		"race_id":     {"db": "race_id", "json": "race"},
	},
	"material_stats": {
		"id":          {"db": "id", "json": ""},
		"material_id": {"db": "material_id", "json": ""},
		"stat_id":     {"db": "stat_id", "json": ""},
		"stat_id2":    {"db": "stat_id2", "json": ""},
	},
	"perks": {
		"id":                {"db": "id", "json": "id"},
		"name":              {"db": "name", "json": "name"},
		"description":       {"db": "description", "json": "description"},
		"icon":              {"db": "icon", "json": "icon"},
		"specialization_id": {"db": "specialization_id", "json": "specialization"},
	},
	"spells": {
		"id":          {"db": "id", "json": "id"},
		"name":        {"db": "name", "json": "name"},
		"description": {"db": "description", "json": "description"},
		"icon":        {"db": "icon", "json": ""},
		"charges":     {"db": "charges", "json": "maxCharges"},
		"class_id":    {"db": "class_id", "json": "class"},
	},
	"spellProperties": {
		"id":          {"db": "id", "json": "id"},
		"name":        {"db": "name", "json": "longDescription"},
		"material_id": {"db": "material_id", "json": "materialId"},
	},
	"stats": {
		"id":        {"db": "id", "json": ""},
		"stat_type": {"db": "stat_type", "json": ""},
	},
}

var correllatedData map[string]map[int]map[string]interface{}

const gameDataRootPath = "/app/gameData"

func ReadJSONFromFileUnstructured(filePath string) ([]map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var payload []map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		slog.Info("Loaded JSON file", "path", filePath, "items", len(payload), "error", err)
		return nil, err
	}

	return payload, nil
}

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
			return classID[0]
		}
		slog.Warn("Invalid class name format", "rawValue", rawValue)
		return nil
	}
	return rawValue
}

func TransformRaceForDB(tableName string, dbFieldName string, rawValue interface{}, result map[string]map[int]map[string]interface{}) interface{} {
	if dbFieldName == "race_id" {
		if raceName, ok := rawValue.(string); ok {
			for id, fields := range result["races"] {
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

func TransformValueForDB(tableName string, dbFieldName string, rawValue interface{}, result map[string]map[int]map[string]interface{}) interface{} {
	rawValue = TransformIconsForDB(tableName, dbFieldName, rawValue)
	rawValue = TransformClassForDB(tableName, dbFieldName, rawValue)
	rawValue = TransformRaceForDB(tableName, dbFieldName, rawValue, result)
	return rawValue
}

func CorrelateFieldsForRow(jsonMeta jsonMeta, rowIndex int, fieldMapping map[string]map[string]string, result map[string]map[int]map[string]interface{}) map[string]interface{} {
	rowData := make(map[string]interface{})

	for _, correlatedField := range fieldMapping {
		jsonFieldName := correlatedField["json"]
		dbFieldName := correlatedField["db"]

		if value, exists := GetFieldValue(jsonMeta, jsonFieldName, rowIndex); exists {
			rowData[dbFieldName] = TransformValueForDB(jsonMeta.Name, dbFieldName, value, result)
		}
	}

	return rowData
}

func ProcessTableRow(jsonMeta jsonMeta, rowIndex int, idValue interface{}, result map[string]map[int]map[string]interface{}) {
	id, ok := GetIDValue(idValue)
	if !ok {
		slog.Warn("Skipping row due to missing or invalid id", "name", jsonMeta.Name, "id", idValue)
		return
	}

	fieldMapping, exists := correllatedFieldNamesMetaMap[jsonMeta.Name]
	if !exists {
		slog.Warn("No field mapping found for table", "name", jsonMeta.Name)
		return
	}

	if result[jsonMeta.Name][id] == nil {
		result[jsonMeta.Name][id] = make(map[string]interface{})
	}

	rowData := CorrelateFieldsForRow(jsonMeta, rowIndex, fieldMapping, result)
	for field, value := range rowData {
		result[jsonMeta.Name][id][field] = value
	}
}

func ProcessTableData(jsonMeta jsonMeta, result map[string]map[int]map[string]interface{}) {
	if result[jsonMeta.Name] == nil {
		result[jsonMeta.Name] = make(map[int]map[string]interface{})
	}

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

func PopulateStaticTable(tableName string, sourceData map[string][2]interface{}, result map[string]map[int]map[string]interface{}) {
	if result[tableName] == nil {
		result[tableName] = make(map[int]map[string]interface{})
	}

	for itemName, itemData := range sourceData {
		id, ok := itemData[0].(int)
		if !ok {
			slog.Warn("Invalid id type in static table", "tableName", tableName, "itemName", itemName, "id", itemData[0])
			continue
		}

		iconPath, ok := itemData[1].(string)
		if !ok {
			slog.Warn("Invalid icon path type in static table", "tableName", tableName, "itemName", itemName, "iconPath", itemData[1])
			continue
		}

		if result[tableName][id] == nil {
			result[tableName][id] = make(map[string]interface{})
		}

		result[tableName][id]["id"] = id
		result[tableName][id]["name"] = itemName

		if iconBytes, ok := IconToBytes(iconPath); ok {
			result[tableName][id]["icon"] = iconBytes
		} else {
			slog.Warn("Failed to read icon for static table", "tableName", tableName, "itemName", itemName, "iconPath", iconPath)
			result[tableName][id]["icon"] = nil
		}
	}

	slog.Info("Populated static table", "tableName", tableName, "count", len(result[tableName]))
}

func PopulateMaterialStats(materialsJsonMeta jsonMeta, result map[string]map[int]map[string]interface{}) {
	if result["material_stats"] == nil {
		result["material_stats"] = make(map[int]map[string]interface{})
	}

	statsFieldData, exists := materialsJsonMeta.Fields["stats"]
	if !exists {
		slog.Warn("No stats field found in materials JSON")
		return
	}

	idFieldData, exists := materialsJsonMeta.Fields["id"]
	if !exists {
		slog.Warn("No id field found in materials JSON")
		return
	}

	idCounter := 0
	for rowIndex, materialIDRaw := range idFieldData {
		materialID, ok := GetIDValue(materialIDRaw)
		if !ok {
			continue
		}

		statNamesRaw, exists := statsFieldData[rowIndex]
		if !exists {
			continue
		}

		statNamesArray, ok := statNamesRaw.([]interface{})
		if !ok {
			slog.Warn("Invalid stats field type for material", "materialID", materialID, "stats", statNamesRaw)
			continue
		}

		if len(statNamesArray) == 0 {
			continue
		}

		materialStatsEntry := map[string]interface{}{
			"id":          idCounter,
			"material_id": materialID,
		}

		if len(statNamesArray) > 0 {
			statName, ok := statNamesArray[0].(string)
			if !ok {
				slog.Warn("Invalid stat name type for material", "materialID", materialID, "statName", statNamesArray[0])
				continue
			}

			statID, found := findStatIDByName(statName, result["stats"])
			if !found {
				slog.Warn("Stat name not found in stats table", "materialID", materialID, "statName", statName)
				continue
			}
			materialStatsEntry["stat_id"] = statID
		}

		if len(statNamesArray) > 1 {
			statName, ok := statNamesArray[1].(string)
			if !ok {
				slog.Warn("Invalid stat name type for material", "materialID", materialID, "statName", statNamesArray[1])
				materialStatsEntry["stat_id2"] = nil
			} else {
				// Look up stat ID by name in the stats table
				statID, found := findStatIDByName(statName, result["stats"])
				if !found {
					slog.Warn("Stat name not found in stats table", "materialID", materialID, "statName", statName)
					materialStatsEntry["stat_id2"] = nil
				} else {
					materialStatsEntry["stat_id2"] = statID
				}
			}
		} else {
			materialStatsEntry["stat_id2"] = nil
		}

		result["material_stats"][idCounter] = materialStatsEntry
		idCounter++
	}

	slog.Info("Populated material_stats table", "count", len(result["material_stats"]))
}

func findStatIDByName(statName string, statsTable map[int]map[string]interface{}) (int, bool) {
	for id, fields := range statsTable {
		if name, ok := fields["name"].(string); ok && name == statName {
			return id, true
		}
	}
	return 0, false
}

func CorrellateJSONDataWithDBFields() (map[string]map[int]map[string]interface{}, error) {
	result := make(map[string]map[int]map[string]interface{})

	var materialsJsonMeta jsonMeta
	for _, jsonMeta := range jsonMetaArray {
		if jsonMeta.Name == "materials" {
			materialsJsonMeta = jsonMeta
		}
		ProcessTableData(jsonMeta, result)
	}

	PopulateStaticTable("classes", classes, result)
	PopulateStaticTable("stats", stats, result)
	PopulateMaterialStats(materialsJsonMeta, result)

	return result, nil
}

func ReadPayload() error {
	for i, jsonMeta := range jsonMetaArray {
		payload, err := ReadJSONFromFileUnstructured(jsonMeta.FilePath)
		if err != nil {
			slog.Error("Error reading payload", "filePath", jsonMeta.FilePath, "error", err)
			return err
		}

		jsonMetaArray[i].Fields = make(map[string]map[int]interface{})
		if len(payload) > 0 {
			for fieldName := range payload[0] {
				jsonMetaArray[i].Fields[fieldName] = make(map[int]interface{})
			}

			for rowIndex, obj := range payload {
				for fieldName, value := range obj {
					jsonMetaArray[i].Fields[fieldName][rowIndex] = value
				}
			}
		}

		slog.Info("Successfully read payload", "filePath", jsonMeta.FilePath, "items", len(payload), "fields", len(jsonMetaArray[i].Fields))
	}
	return nil
}

func Run() {
	readErr := ReadPayload()
	if readErr != nil {
		slog.Error("Error reading payload", "error", readErr)
	}

	var table, err = CorrellateJSONDataWithDBFields()
	if err != nil {
		slog.Error("Error correlating JSON data with DB fields", "error", err)
	}

	correllatedData = table

	// if artifactData, exists := correllatedData["material_stats"]; exists && len(artifactData) > 0 {
	// 	for id, fields := range artifactData {
	// 		for field, value := range fields {
	// 			slog.Info("Fields for table", "id", id, "field", field, "value", value)
	// 		}
	// 	}
	// }
}
