package dbseeder

import "log/slog"

type config struct {
	logger                      *slog.Logger
	gameDataRootPath            string
	jsonSources                 map[string]jsonMeta
	staticTables                map[string][]map[string]interface{}
	correlatedFieldNamesMetaMap map[string][]fieldMapping
	arrayToJunctionTableSpecs   map[string]junctionTableSpec
}

const gameDataRootPath = "/app/gameData/"

func newSeederConfig() *config {
	return &config{
		logger:           slog.Default(),
		gameDataRootPath: gameDataRootPath,
		jsonSources: map[string]jsonMeta{
			"artifacts":            {filePath: gameDataRootPath + "artifacts.json", name: "artifacts", isDerived: false},
			"races":                {filePath: gameDataRootPath + "races.json", name: "races", isDerived: false},
			"specializations":      {filePath: gameDataRootPath + "specializations.json", name: "specializations", isDerived: false},
			"materials":            {filePath: gameDataRootPath + "materials.json", name: "materials", isDerived: false},
			"traits":               {filePath: gameDataRootPath + "traits.json", name: "traits", isDerived: false},
			"creatures":            {filePath: gameDataRootPath + "creatures.json", name: "creatures", isDerived: false},
			"perks":                {filePath: gameDataRootPath + "perks.json", name: "perks", isDerived: false},
			"spells":               {filePath: gameDataRootPath + "spells.json", name: "spells", isDerived: false},
			"spellProperties":      {filePath: gameDataRootPath + "spellProperties.json", name: "spellProperties", isDerived: false},
			"relics":               {filePath: gameDataRootPath + "gods.json", name: "relics", isDerived: false},
			"material_stats":       {name: "material_stats", isDerived: true},
			"creature_stat_growth": {name: "creature_stat_growth", isDerived: true},
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
				{dbField: "id", jsonField: "id"},
				{dbField: "name", jsonField: "name"},
				{dbField: "icon", jsonField: "icon"},
			},
			"stats": {
				{dbField: "id", jsonField: "id"},
				{dbField: "type", jsonField: "name"},
				{dbField: "icon", jsonField: "icon"},
			},
			"artifacts": {
				{dbField: "id", jsonField: "id"},
				{dbField: "name", jsonField: "name"},
				{dbField: "icon", jsonField: "icons"},
				{dbField: "stat_id", jsonField: "stat", findIdFromSource: "stats"},
			},
			"races": {
				{dbField: "id", jsonField: "id"},
				{dbField: "name", jsonField: "name"},
				{dbField: "icon", jsonField: "icon"},
			},
			"specializations": {
				{dbField: "id", jsonField: "id"},
				{dbField: "name", jsonField: "name"},
				{dbField: "description", jsonField: "description"},
				{dbField: "icon", jsonField: "icon"},
			},
			"materials": {
				{dbField: "id", jsonField: "id"},
				{dbField: "name", jsonField: "name"},
				{dbField: "description", jsonField: "description"},
				{dbField: "icon", jsonField: "icon"},
				{dbField: "type", jsonField: "slot"},
				{dbField: "", jsonField: "stats"},
			},
			"traits": {
				{dbField: "id", jsonField: "id"},
				{dbField: "name", jsonField: "name"},
				{dbField: "description", jsonField: "description"},
				{dbField: "material_id", jsonField: "item"},
			},
			"perks": {
				{dbField: "id", jsonField: "id"},
				{dbField: "name", jsonField: "name"},
				{dbField: "description", jsonField: "description"},
				{dbField: "icon", jsonField: "icon"},
				{dbField: "specialization_id", jsonField: "specialization", findIdFromSource: "specializations"},
			},
			"spells": {
				{dbField: "id", jsonField: "id"},
				{dbField: "name", jsonField: "name"},
				{dbField: "description", jsonField: "description"},
				{dbField: "charges", jsonField: "maxCharges"},
				{dbField: "class_id", jsonField: "class", findIdFromSource: "classes"},
			},
			"spellProperties": {
				{dbField: "id", jsonField: "id"},
				{dbField: "name", jsonField: "shortDescription"},
				{dbField: "material_id", jsonField: "item"},
			},
			"creatures": {
				{dbField: "id", jsonField: "id"},
				{dbField: "name", jsonField: "name"},
				{dbField: "description", jsonField: "description"},
				{dbField: "icon", jsonField: "battleSprite"},
				{dbField: "", jsonField: "statGrowth"},
				{dbField: "race_id", jsonField: "race", findIdFromSource: "races"},
				{dbField: "class_id", jsonField: "class", findIdFromSource: "classes"},
				{dbField: "trait_id", jsonField: "trait", findIdFromSource: "traits"},
			},
			"relics": {
				{dbField: "id", jsonField: "id"},
				{dbField: "name", jsonField: "relicTitle"},
				{dbField: "bonuses", jsonField: "relicBonuses"},
				{dbField: "icon", jsonField: "relicBigIcon"},
				{dbField: "stat_id", jsonField: "relicStat", findIdFromSource: "stats"},
			},
		},
		arrayToJunctionTableSpecs: map[string]junctionTableSpec{
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
	filePath  string
	name      string
	isDerived bool
	items     []map[string]interface{}
}

type fieldMapping struct {
	dbField          string
	jsonField        string
	findIdFromSource string
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
