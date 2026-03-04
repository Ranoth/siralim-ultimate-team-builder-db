package dbseeder

import "log/slog"

type config struct {
	logger                      *slog.Logger
	gameDataRootPath            string
	jsonSources                 map[string]jsonMeta
	staticTables                map[string][]map[string]interface{}
	correlatedFieldNamesMetaMap map[string][]fieldMapping
	junctionTableSpecs          map[string]junctionTableSpec
}

const gameDataRootPath = "/app/gameData/"

func newSeederConfig() *config {
	return &config{
		logger:           slog.Default(),
		gameDataRootPath: gameDataRootPath,
		jsonSources: map[string]jsonMeta{
			"artifacts":       {filePath: gameDataRootPath + "artifacts.json", name: "artifacts"},
			"races":           {filePath: gameDataRootPath + "races.json", name: "races"},
			"specializations": {filePath: gameDataRootPath + "specializations.json", name: "specializations"},
			"materials":       {filePath: gameDataRootPath + "materials.json", name: "materials"},
			"traits":          {filePath: gameDataRootPath + "traits.json", name: "traits"},
			"creatures":       {filePath: gameDataRootPath + "creatures.json", name: "creatures"},
			"perks":           {filePath: gameDataRootPath + "perks.json", name: "perks"},
			"spells":          {filePath: gameDataRootPath + "spells.json", name: "spells"},
			"spellProperties": {filePath: gameDataRootPath + "spellProperties.json", name: "spellProperties"},
			"relics":          {filePath: gameDataRootPath + "gods.json", name: "relics"},
			"material_stats":  {name: "material_stats"},
		},
		staticTables: map[string][]map[string]interface{}{
			"classes": {
				{"name": "life", "id": 0, "icon": "images/misc/class/life.png"},
				{"name": "death", "id": 1, "icon": "images/misc/class/death.png"},
				{"name": "nature", "id": 2, "icon": "images/misc/class/nature.png"},
				{"name": "sorcery", "id": 3, "icon": "images/misc/class/sorcery.png"},
				{"name": "chaos", "id": 4, "icon": "images/misc/class/chaos.png"},
			},
			"stats": {
				{"name": "health", "id": 0, "icon": "images/misc/stat/health.png"},
				{"name": "attack", "id": 1, "icon": "images/misc/stat/attack.png"},
				{"name": "intelligence", "id": 2, "icon": "images/misc/stat/intelligence.png"},
				{"name": "defense", "id": 3, "icon": "images/misc/stat/defense.png"},
				{"name": "speed", "id": 4, "icon": "images/misc/stat/speed.png"},
			},
		},
		correlatedFieldNamesMetaMap: map[string][]fieldMapping{
			"classes": {
				{dBField: "id", jsonField: "id"},
				{dBField: "name", jsonField: "name"},
				{dBField: "icon", jsonField: "icon"},
			},
			"stats": {
				{dBField: "id", jsonField: "id"},
				{dBField: "type", jsonField: "name"},
				{dBField: "icon", jsonField: "icon"},
			},
			"artifacts": {
				{dBField: "id", jsonField: "id"},
				{dBField: "name", jsonField: "name"},
				{dBField: "icon", jsonField: "icons"},
				{dBField: "stat_id", jsonField: "stat", findIdFromSource: "stats"},
			},
			"races": {
				{dBField: "id", jsonField: "id"},
				{dBField: "name", jsonField: "name"},
				{dBField: "icon", jsonField: "icon"},
			},
			"specializations": {
				{dBField: "id", jsonField: "id"},
				{dBField: "name", jsonField: "name"},
				{dBField: "description", jsonField: "description"},
				{dBField: "icon", jsonField: "icon"},
			},
			"materials": {
				{dBField: "id", jsonField: "id"},
				{dBField: "name", jsonField: "name"},
				{dBField: "description", jsonField: "description"},
				{dBField: "icon", jsonField: "icon"},
				{dBField: "type", jsonField: "slot"},
				{dBField: "", jsonField: "stats"},
			},
			"traits": {
				{dBField: "id", jsonField: "id"},
				{dBField: "name", jsonField: "name"},
				{dBField: "description", jsonField: "description"},
				{dBField: "material_id", jsonField: "item"},
			},
			"perks": {
				{dBField: "id", jsonField: "id"},
				{dBField: "name", jsonField: "name"},
				{dBField: "description", jsonField: "description"},
				{dBField: "icon", jsonField: "icon"},
				{dBField: "specialization_id", jsonField: "specialization", findIdFromSource: "specializations"},
			},
			"spells": {
				{dBField: "id", jsonField: "id"},
				{dBField: "name", jsonField: "name"},
				{dBField: "description", jsonField: "description"},
				{dBField: "charges", jsonField: "maxCharges"},
				{dBField: "class_id", jsonField: "class", findIdFromSource: "classes"},
			},
			"spellProperties": {
				{dBField: "id", jsonField: "id"},
				{dBField: "name", jsonField: "shortDescription"},
				{dBField: "material_id", jsonField: "item"},
			},
			"creatures": {
				{dBField: "id", jsonField: "id"},
				{dBField: "name", jsonField: "name"},
				{dBField: "description", jsonField: "description"},
				{dBField: "icon", jsonField: "battleSprite"},
				{dBField: "race_id", jsonField: "race", findIdFromSource: "races"},
				{dBField: "class_id", jsonField: "class", findIdFromSource: "classes"},
				{dBField: "trait_id", jsonField: "trait", findIdFromSource: "traits"},
			},
			"relics": {
				{dBField: "id", jsonField: "id"},
				{dBField: "name", jsonField: "relicTitle"},
				{dBField: "bonuses", jsonField: "relicBonuses"},
				{dBField: "icon", jsonField: "relicBigIcon"},
				{dBField: "stat_id", jsonField: "relicStat", findIdFromSource: "stats"},
			},
		},
		junctionTableSpecs: map[string]junctionTableSpec{
			"material_stats": {
				name:           "material_stats",
				sourceName:     "materials",
				dataField:      "stats",
				parentKeyField: "id",
				mappings: []junctionFieldMapping{
					{junctionField: "material_id", sourceField: "id", arrayIndex: -1},
					{junctionField: "stat_id", sourceField: "stats", arrayIndex: 0, findIdFromSource: "stats"},
					{junctionField: "stat_id2", sourceField: "stats", arrayIndex: 1, findIdFromSource: "stats"},
				},
			},
		},
	}
}

type jsonMeta struct {
	filePath string
	name     string
	items    []map[string]interface{}
}

type fieldMapping struct {
	dBField          string
	jsonField        string
	findIdFromSource string
	isFake           bool
}

type junctionTableSpec struct {
	name           string
	sourceName     string
	dataField      string
	parentKeyField string // The field in the source item that contains the parent ID (e.g., "id" for materials)
	mappings       []junctionFieldMapping
}

type junctionFieldMapping struct {
	junctionField    string
	sourceField      string
	arrayIndex       int // -1 if not an array element, 0 for first, 1 for second, etc.
	findIdFromSource string
}
