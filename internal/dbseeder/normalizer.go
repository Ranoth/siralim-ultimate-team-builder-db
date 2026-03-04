package dbseeder

import (
	"log/slog"
	"os"
	"path/filepath"
)

type normalizer struct {
	jsonParser *jsonParser
	logger     *slog.Logger
	config     *config
}

func newNormalizer(jsonParser *jsonParser, logger *slog.Logger, config *config) *normalizer {
	return &normalizer{jsonParser: jsonParser, logger: logger, config: config}
}

func (t *normalizer) iconToBytes(iconPathValue interface{}) ([]byte, bool) {
	readIcon := func(iconPath string) ([]byte, bool) {
		absolutePath := filepath.Join(t.config.gameDataRootPath, iconPath)
		iconBytes, err := os.ReadFile(absolutePath)
		if err != nil {
			t.logger.Warn("Skipping icon conversion", "iconPath", absolutePath, "error", err)
			return nil, false
		}

		return iconBytes, true
	}

	switch value := iconPathValue.(type) {
	case string:
		if value == "" {
			return nil, false
		}
		return readIcon(value)
	case []string:
		for _, iconPath := range value {
			if iconPath == "" {
				continue
			}
			if iconBytes, ok := readIcon(iconPath); ok {
				return iconBytes, true
			}
		}
	case []interface{}:
		for _, rawPath := range value {
			iconPath, ok := rawPath.(string)
			if !ok || iconPath == "" {
				continue
			}
			if iconBytes, ok := readIcon(iconPath); ok {
				return iconBytes, true
			}
		}
	}

	t.logger.Warn("Skipping icon conversion", "reason", "unsupported or empty iconPath type", "type", iconPathValue)
	return nil, false
}

func (t *normalizer) replaceLocalValueWithForeignValue(valueNameOne string, valueNameTwo string) {
	cfnmmap := t.config.correlatedFieldNamesMetaMap
	sources := t.config.jsonSources
	referenceIndexes := make(map[string]map[interface{}]interface{})
	referenceIDSets := make(map[string]map[interface{}]struct{})

	for sourceName, sourceData := range sources {
		index := make(map[interface{}]interface{}, len(sourceData.items))
		idSet := make(map[interface{}]struct{}, len(sourceData.items))
		for _, item := range sourceData.items {
			valueOne, hasValueOne := item[valueNameOne]
			valueTwo, hasValueTwo := item[valueNameTwo]

			if !hasValueOne || !hasValueTwo {
				continue
			}
			index[valueOne] = valueTwo
			idSet[valueTwo] = struct{}{}
		}
		referenceIndexes[sourceName] = index
		referenceIDSets[sourceName] = idSet
	}

	for sourceName, fieldMappings := range cfnmmap {
		sourceData := sources[sourceName]

		for _, fieldMapping := range fieldMappings {
			if fieldMapping.findIdFromSource == "" {
				continue
			}

			referenceIndex := referenceIndexes[fieldMapping.findIdFromSource]
			referenceIDSet := referenceIDSets[fieldMapping.findIdFromSource]
			if len(referenceIndex) == 0 {
				continue
			}

			for i := range sourceData.items {
				jsonValue := sourceData.items[i][fieldMapping.jsonField]
				if jsonValue == nil {
					continue
				}

				if idValue, ok := referenceIndex[jsonValue]; ok {
					sourceData.items[i][fieldMapping.jsonField] = idValue
					continue
				}

				// Some JSON fields already contain direct foreign-key IDs.
				if _, ok := referenceIDSet[jsonValue]; ok {
					continue
				}

				sourceData.items[i][fieldMapping.jsonField] = nil
			}
		}

		sources[sourceName] = sourceData
	}
}

func (t *normalizer) removeInvalidDirectForeignKeyIDs() {
	sources := t.config.jsonSources

	refIDs := make(map[string]map[interface{}]struct{})
	for sourceName, sourceData := range sources {
		ids := make(map[interface{}]struct{}, len(sourceData.items))
		for _, item := range sourceData.items {
			if idValue, ok := item["id"]; ok {
				ids[idValue] = struct{}{}
			}
		}
		refIDs[sourceName] = ids
	}

	validators := []struct {
		sourceName string
		fieldName  string
		refSource  string
	}{
		{sourceName: "traits", fieldName: "item", refSource: "materials"},
		{sourceName: "spellProperties", fieldName: "item", refSource: "materials"},
	}

	for _, validator := range validators {
		sourceData, sourceOK := sources[validator.sourceName]
		validIDs, refOK := refIDs[validator.refSource]
		if !sourceOK || !refOK || len(validIDs) == 0 {
			continue
		}

		for i := range sourceData.items {
			rawValue, exists := sourceData.items[i][validator.fieldName]
			if !exists || rawValue == nil {
				continue
			}

			if _, ok := validIDs[rawValue]; !ok {
				sourceData.items[i][validator.fieldName] = nil
				t.logger.Warn("Removed invalid foreign key id",
					"source", validator.sourceName,
					"field", validator.fieldName,
					"value", rawValue,
					"referenceSource", validator.refSource)
			}
		}

		sources[validator.sourceName] = sourceData
	}
}

func (t *normalizer) renameFieldsToDbNames() {
	cfnmmap := t.config.correlatedFieldNamesMetaMap
	sources := t.config.jsonSources

	for sourceName, fieldMappings := range cfnmmap {
		sourceData := sources[sourceName]

		for i := range sourceData.items {
			for _, fieldMapping := range fieldMappings {
				if fieldMapping.jsonField == "" || fieldMapping.jsonField == fieldMapping.dBField {
					continue
				}

				if value, exists := sourceData.items[i][fieldMapping.jsonField]; exists {
					sourceData.items[i][fieldMapping.dBField] = value
					delete(sourceData.items[i], fieldMapping.jsonField)
				}
			}
		}

		sources[sourceName] = sourceData
	}
}

func (t *normalizer) convertIconPathsToBytes() {
	cfnmmap := t.config.correlatedFieldNamesMetaMap
	sources := t.config.jsonSources
	nullableIconSources := map[string]struct{}{
		"races":           {},
		"specializations": {},
		"relics":          {},
		"creatures":       {},
	}

	for sourceName, fieldMappings := range cfnmmap {
		sourceData := sources[sourceName]

		// Find icon fields for this source
		iconFields := make([]string, 0)
		for _, fieldMapping := range fieldMappings {
			if fieldMapping.dBField == "icon" {
				iconFields = append(iconFields, fieldMapping.dBField)
			}
		}

		if len(iconFields) == 0 {
			continue
		}

		// Convert icon paths to bytes for each item
		for i := range sourceData.items {
			for _, iconField := range iconFields {
				iconValue := sourceData.items[i][iconField]
				if iconValue != nil {
					if iconBytes, ok := t.iconToBytes(iconValue); ok {
						sourceData.items[i][iconField] = iconBytes
					} else {
						if _, isNullableIconSource := nullableIconSources[sourceName]; isNullableIconSource {
							sourceData.items[i][iconField] = nil
						}
					}
				}
			}
		}

		sources[sourceName] = sourceData
	}
}

func (t *normalizer) convertArrayDescriptionsToStrings() {
	cfnmmap := t.config.correlatedFieldNamesMetaMap
	sources := t.config.jsonSources

	for sourceName, fieldMappings := range cfnmmap {
		sourceData := sources[sourceName]

		// Find description fields for this source
		descFields := make([]string, 0)
		for _, fieldMapping := range fieldMappings {
			if fieldMapping.dBField == "description" {
				descFields = append(descFields, fieldMapping.dBField)
			}
		}

		if len(descFields) == 0 {
			continue
		}

		// Convert array descriptions to strings for each item
		for i := range sourceData.items {
			for _, descField := range descFields {
				descValue := sourceData.items[i][descField]
				if descValue != nil {
					switch v := descValue.(type) {
					case []interface{}:
						if len(v) > 0 {
							if str, ok := v[0].(string); ok {
								sourceData.items[i][descField] = str
							}
						}
					case []string:
						if len(v) > 0 {
							sourceData.items[i][descField] = v[0]
						}
					}
				}
			}
		}

		sources[sourceName] = sourceData
	}
}

func (t *normalizer) seedJunctionTables() {
	sources := t.config.jsonSources

	for junctionName, spec := range t.config.junctionTableSpecs {
		sourceData := sources[spec.sourceName]
		junctionItems := make([]map[string]interface{}, 0)
		nextID := 0

		referenceIndexes := make(map[string]map[interface{}]interface{})
		for _, mapping := range spec.mappings {
			if mapping.findIdFromSource == "" {
				continue
			}

			refSourceData := sources[mapping.findIdFromSource]
			refIndex := make(map[interface{}]interface{}, len(refSourceData.items))
			for _, item := range refSourceData.items {
				if nameValue, ok := item["name"]; ok {
					if idValue, ok := item["id"]; ok {
						refIndex[nameValue] = idValue
					}
				}
			}
			referenceIndexes[mapping.findIdFromSource] = refIndex
		}

		for _, sourceItem := range sourceData.items {
			arrayValue, ok := sourceItem[spec.dataField]
			if !ok {
				continue
			}

			var arrayItems []string
			switch v := arrayValue.(type) {
			case []interface{}:
				for _, item := range v {
					if str, ok := item.(string); ok {
						arrayItems = append(arrayItems, str)
					}
				}
			}

			if len(arrayItems) == 0 {
				continue
			}

			junctionRecord := make(map[string]interface{})
			junctionRecord["id"] = nextID
			nextID++

			for _, mapping := range spec.mappings {
				var fieldValue interface{}

				if mapping.arrayIndex == -1 {
					fieldValue = sourceItem[mapping.sourceField]
				} else if mapping.arrayIndex < len(arrayItems) {
					fieldValue = arrayItems[mapping.arrayIndex]
				} else {
					continue
				}

				if fieldValue == nil {
					continue
				}

				if mapping.findIdFromSource != "" {
					if refIndex, ok := referenceIndexes[mapping.findIdFromSource]; ok {
						if idValue, ok := refIndex[fieldValue]; ok {
							junctionRecord[mapping.junctionField] = idValue
						}
					}
				} else {
					junctionRecord[mapping.junctionField] = fieldValue
				}
			}

			junctionItems = append(junctionItems, junctionRecord)
		}

		sources[junctionName] = jsonMeta{
			name:  junctionName,
			items: junctionItems,
		}

		t.logger.Info("Seeded junction table", "name", junctionName, "records", len(junctionItems))
	}
}

func (t *normalizer) removeNullAndEmptyFields() {
	sources := t.config.jsonSources

	for sourceName, sourceData := range sources {
		for i := range sourceData.items {
			fieldsToDelete := make([]string, 0)
			for fieldName, fieldValue := range sourceData.items[i] {
				// Mark field for deletion if it's nil
				if fieldValue == nil {
					fieldsToDelete = append(fieldsToDelete, fieldName)
					continue
				}

				// Mark field for deletion if it's an empty array
				switch v := fieldValue.(type) {
				case []interface{}:
					if len(v) == 0 {
						fieldsToDelete = append(fieldsToDelete, fieldName)
					}
				case []string:
					if len(v) == 0 {
						fieldsToDelete = append(fieldsToDelete, fieldName)
					}
				}
			}

			// Delete marked fields
			for _, fieldName := range fieldsToDelete {
				delete(sourceData.items[i], fieldName)
			}
		}

		sources[sourceName] = sourceData
	}
}

func (t *normalizer) normalize() {
	t.replaceLocalValueWithForeignValue("name", "id")
	t.removeInvalidDirectForeignKeyIDs()
	t.seedJunctionTables()
	t.renameFieldsToDbNames()
	t.convertArrayDescriptionsToStrings()
	t.convertIconPathsToBytes()
	t.removeNullAndEmptyFields()

	for sourceName, jsonMeta := range t.config.jsonSources {
		if len(jsonMeta.items) > 0 {
			if (sourceName == "traits") && len(jsonMeta.items) > 1 {
				t.logger.Info("Sample after normalization",
					"source", sourceName,
					"item", jsonMeta.items[675])
			}
		}
	}
}
