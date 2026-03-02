package dbseeder

import (
	"encoding/json"
	"log/slog"
	"os"
)

func readJSONFromFileUnstructured(filePath string) ([]map[string]interface{}, error) {
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

func correlateJSONDataWithDBFields() (*TransformedData, error) {
	result := NewTransformedData()

	// Tables that use explicit mappers (skip generic processing)
	explicitMappers := map[string]bool{
		"creatures":      true,
		"material_stats": true, // Derived from materials.json
	}

	var materialsJsonMeta jsonMeta
	var creaturesJsonMeta jsonMeta

	for _, jsonMeta := range jsonSources {
		if jsonMeta.Name == "materials" {
			materialsJsonMeta = jsonMeta
		}
		if jsonMeta.Name == "creatures" {
			creaturesJsonMeta = jsonMeta
		}

		// Skip tables with explicit mappers
		if explicitMappers[jsonMeta.Name] {
			continue
		}

		ProcessTableData(jsonMeta, result)
	}

	// Populate static tables first (needed for lookups)
	PopulateStaticTable("classes", classes, result)
	PopulateStaticTable("stats", stats, result)

	// Process tables with explicit mappers (after dependencies are ready)
	ProcessMaterialStatsExplicit(materialsJsonMeta.RawPayload, result)
	ProcessCreaturesExplicit(creaturesJsonMeta.RawPayload, result)

	return result, nil
}

func loadPayload() error {
	for i, jsonMeta := range jsonSources {
		payload, err := readJSONFromFileUnstructured(jsonMeta.FilePath)
		if err != nil {
			slog.Error("Error reading payload", "filePath", jsonMeta.FilePath, "error", err)
			return err
		}

		// Store raw payload for explicit mappers
		jsonSources[i].RawPayload = payload

		jsonSources[i].Fields = make(map[string]map[int]interface{})
		if len(payload) > 0 {
			for fieldName := range payload[0] {
				jsonSources[i].Fields[fieldName] = make(map[int]interface{})
			}

			for rowIndex, obj := range payload {
				for fieldName, value := range obj {
					jsonSources[i].Fields[fieldName][rowIndex] = value
				}
			}
		}

		slog.Info("Successfully read payload", "filePath", jsonMeta.FilePath, "items", len(payload), "fields", len(jsonSources[i].Fields))
	}
	return nil
}

func Run() {
	readErr := loadPayload()
	if readErr != nil {
		slog.Error("Error reading payload", "error", readErr)
	}

	var table, err = correlateJSONDataWithDBFields()
	if err != nil {
		slog.Error("Error correlating JSON data with DB fields", "error", err)
	}

	correlatedTables = table
}

// GetCorrelatedTables returns the transformed data ready for database insertion
func GetCorrelatedTables() *TransformedData {
	return correlatedTables
}
