package dbseeder

import (
	"encoding/json"
	"log/slog"
	"os"
)

var tablesFilePathsAbsoluteMap = []map[string]interface{}{
	{"filePath": "/app/gameData/artifacts.json", "tables": map[string][]string{"artifacts": {"id", "name", "description", "icon"}}},
	{"filePath": "/app/gameData/creatures.json", "tables": map[string][]string{"creatures": {"id", "name", "description", "icon", "trait_id", "class_id", "race_id"}}},
	{"filePath": "/app/gameData/class.json", "tables": map[string][]string{"class": {"id", "name", "icon"}}},
	{"filePath": "/app/gameData/materials.json", "tables": map[string][]string{"materials": {"id", "name", "description", "icon", "type"}}},
	{"filePath": "virtual", "tables": map[string][]string{"material_stats": {"id", "material_id", "stat_id", "value"}}},
	{"filePath": "/app/gameData/perks.json", "tables": map[string][]string{"perks": {"id", "name", "description", "icon", "specialization_id"}}},
	{"filePath": "/app/gameData/races.json", "tables": map[string][]string{"races": {"id", "name", "icon"}}},
	{"filePath": "/app/gameData/specializations.json", "tables": map[string][]string{"specializations": {"id", "name", "description"}}},
	{"filePath": "/app/gameData/spells.json", "tables": map[string][]string{"spells": {"id", "name", "description", "icon", "charges", "class_id"}}},
	{"filePath": "/app/gameData/spellProperties.json", "tables": map[string][]string{"spellProperties": {"id", "name", "description", "material_id"}}},
	{"filePath": "virtual", "tables": map[string][]string{"stats": {"id", "type"}}},
	{"filePath": "/app/gameData/traits.json", "tables": map[string][]string{"traits": {"id", "name", "description", "material_id"}}},
}

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

func ReadPayload() error {
	for _, filePath := range tablesFilePathsAbsoluteMap {
		payload, err := ReadJSONFromFileUnstructured(filePath["filePath"].(string))
		if err != nil {
			slog.Error("Error reading payload", "filePath", filePath, "error", err)
			return err
		}
		slog.Info("Successfully read payload", "filePath", filePath, "items", len(payload))

	}
	return nil
}

func Run() {
	err := ReadPayload()
	if err != nil {
		slog.Error("Error reading payload", "error", err)
	}
}
