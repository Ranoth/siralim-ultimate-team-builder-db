package dbseeder

import (
	"encoding/json"
	"log/slog"
	"os"
)

type jsonParser struct {
	logger *slog.Logger
	config *config
}

func newJSONParser(logger *slog.Logger, config *config) *jsonParser {
	return &jsonParser{logger: logger, config: config}
}

func (p *jsonParser) readJSONFromFileUnStructured(filePath string) ([]map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var payload []map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		p.logger.Info("Loaded JSON file", "path", filePath, "items", len(payload), "error", err)
		return nil, err
	}

	return payload, nil
}

func (p *jsonParser) insertStaticDataToSources() {
	for tableName, tableData := range p.config.staticTables {
		meta := p.config.jsonSources[tableName]
		meta.items = append(meta.items, tableData...)
		p.config.jsonSources[tableName] = meta
		p.logger.Info("Inserted static data into sources", "table", tableName, "items", len(tableData))
	}
}

func (p *jsonParser) maxIDInPayload(payload []map[string]interface{}) int {
	maxID := -1

	for i := range payload {
		if v, ok := payload[i]["id"]; ok {
			if f, ok := v.(float64); ok {
				id := int(f)
				if id > maxID {
					maxID = id
				}
			}
		}
	}

	return maxID
}

func (p *jsonParser) attributeMissingIDsAndCull(sourceName string, payload []map[string]interface{}) []map[string]interface{} {
	fieldMappings, ok := p.config.correlatedFieldNamesMetaMap[sourceName]
	if !ok || len(payload) == 0 {
		return payload
	}

	// Track which json fields to keep
	allowedJsonFields := make(map[string]struct{})
	keepID := false

	for _, fieldMapping := range fieldMappings {
		if fieldMapping.jsonField != "" {
			allowedJsonFields[fieldMapping.jsonField] = struct{}{}
		}

		if fieldMapping.dBField == "id" {
			keepID = true
		}
	}

	nextID := p.maxIDInPayload(payload) + 1
	transformed := make([]map[string]interface{}, len(payload))

	for i, item := range payload {
		filtered := make(map[string]interface{})

		idValue, hasID := item["id"]
		if !hasID {
			idValue = nextID
			nextID++
		}

		if keepID {
			filtered["id"] = idValue
		}

		for jsonFieldName := range allowedJsonFields {
			if value, exists := item[jsonFieldName]; exists {
				filtered[jsonFieldName] = value
			}
		}

		transformed[i] = filtered
	}

	return transformed
}

func (p *jsonParser) parseAndStore() {
	for key, value := range p.config.jsonSources {
		payload, err := p.readJSONFromFileUnStructured(value.filePath)
		if err != nil {
			p.logger.Info("Failed to read JSON file", "path", value.filePath, "error", err)
			continue
		}

		payload = p.attributeMissingIDsAndCull(key, payload)

		source := p.config.jsonSources[key]
		source.items = payload
		p.config.jsonSources[key] = source

		p.logger.Info("Loaded JSON file into sources", "path", p.config.jsonSources[key].filePath, "items", len(payload))
	}

	p.insertStaticDataToSources()
}
