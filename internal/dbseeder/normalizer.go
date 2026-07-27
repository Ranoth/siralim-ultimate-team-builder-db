package dbseeder

import (
	"errors"
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

func (t *normalizer) buildNameToIDIndex(sourceName string) map[interface{}]interface{} {
	index := make(map[interface{}]interface{})
	source, ok := t.config.jsonSources[sourceName]
	if !ok {
		return index
	}

	for _, item := range source.items {
		nameValue, hasName := item["name"]
		idValue, hasID := item["id"]
		if hasName && hasID {
			index[nameValue] = idValue
		}
	}

	return index
}

func (t *normalizer) buildFieldValueSet(sourceName string, fieldName string) map[interface{}]struct{} {
	set := make(map[interface{}]struct{})
	source, ok := t.config.jsonSources[sourceName]
	if !ok {
		return set
	}

	for _, item := range source.items {
		if value, exists := item[fieldName]; exists {
			set[value] = struct{}{}
		}
	}

	return set
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
		validIDs := t.buildFieldValueSet(validator.refSource, "id")
		_, refOK := sources[validator.refSource]
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
				if fieldMapping.jsonField == "" || fieldMapping.jsonField == fieldMapping.dbField {
					continue
				}

				if fieldMapping.dbField == "" {
					continue
				}

				if value, exists := sourceData.items[i][fieldMapping.jsonField]; exists {
					sourceData.items[i][fieldMapping.dbField] = value
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
			if fieldMapping.dbField == "icon" {
				iconFields = append(iconFields, fieldMapping.dbField)
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
			if fieldMapping.dbField == "description" {
				descFields = append(descFields, fieldMapping.dbField)
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

func (t *normalizer) seedArraysToJunctionTables() error {
	sources := t.config.jsonSources

	for junctionName, spec := range t.config.arrayJunctionTableSpecs {
		sourceData := sources[spec.sourceName]
		junctionItems := make([]map[string]interface{}, 0)
		nextID := 0

		referenceIndexes := make(map[string]map[interface{}]interface{})
		for _, mapping := range spec.mappings {
			if mapping.findIdFromSource == "" {
				continue
			}

			if _, exists := referenceIndexes[mapping.findIdFromSource]; !exists {
				referenceIndexes[mapping.findIdFromSource] = t.buildNameToIDIndex(mapping.findIdFromSource)
			}
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

		if len(junctionItems) == 0 {
			return errors.New("no records found for junction table: " + junctionName)
		}

		t.logger.Info("Seeded junction table", "name", junctionName, "records", len(junctionItems))
	}

	return nil
}

func (t *normalizer) seedObjectsToJunctionTables() error {
	sources := t.config.jsonSources

	for junctionName, spec := range t.config.objectJunctionTableSpecs {
		sourceData, ok := sources[spec.sourceName]
		if !ok {
			return errors.New("source not found for junction table: " + junctionName)
		}

		referenceIndexes := make(map[string]map[interface{}]interface{})
		for _, mapping := range spec.mappings {
			if mapping.findIdFromSource == "" {
				continue
			}
			if _, exists := referenceIndexes[mapping.findIdFromSource]; !exists {
				referenceIndexes[mapping.findIdFromSource] = t.buildNameToIDIndex(mapping.findIdFromSource)
			}
		}

		junctionItems := make([]map[string]interface{}, 0)
		nextID := 0

		for _, sourceItem := range sourceData.items {
			objectRaw, exists := sourceItem[spec.dataField]
			if !exists || objectRaw == nil {
				continue
			}

			objectMap, ok := objectRaw.(map[string]interface{})
			if !ok {
				continue
			}

			for objectKey, objectValue := range objectMap {
				record := map[string]interface{}{
					"id": nextID,
				}
				nextID++

				for _, mapping := range spec.mappings {
					switch mapping.sourceKind {
					case "parent":
						parentValue, hasParent := sourceItem[mapping.sourceField]
						if !hasParent || parentValue == nil {
							continue
						}
						record[mapping.junctionField] = parentValue

					case "objectKey":
						value := interface{}(objectKey)

						if mapping.findIdFromSource != "" {
							refIndex := referenceIndexes[mapping.findIdFromSource]
							if refID, ok := refIndex[value]; ok {
								value = refID
							} else {
								value = nil
							}
						}

						if value != nil {
							record[mapping.junctionField] = value
						}

					case "objectValue":
						value := objectValue

						record[mapping.junctionField] = value
					}
				}

				complete := true
				for _, mapping := range spec.mappings {
					if _, ok := record[mapping.junctionField]; !ok {
						complete = false
						break
					}
				}
				if !complete {
					continue
				}

				junctionItems = append(junctionItems, record)
			}
		}

		sources[junctionName] = jsonMeta{
			name:  junctionName,
			items: junctionItems,
		}

		if len(junctionItems) == 0 {
			return errors.New("no records found for junction table: " + junctionName)
		}

		t.logger.Info("Seeded junction table", "name", junctionName, "records", len(junctionItems))
	}

	return nil
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

func (t *normalizer) normalize() error {
	// ORDER MATTERS HERE
	// Each step depends on the data being in a certain state
	t.replaceLocalValueWithForeignValue("name", "id")
	if err := t.seedArraysToJunctionTables(); err != nil {
		return err
	}
	if err := t.seedObjectsToJunctionTables(); err != nil {
		return err
	}
	t.removeInvalidDirectForeignKeyIDs()
	t.renameFieldsToDbNames()
	t.convertArrayDescriptionsToStrings()
	t.convertIconPathsToBytes()
	t.removeNullAndEmptyFields()

	return nil
}
