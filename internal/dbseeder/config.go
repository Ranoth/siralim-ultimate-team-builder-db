package dbseeder

type jsonMeta struct {
	FilePath   string
	Name       string
	Fields     map[string]map[int]interface{}
	RawPayload []map[string]interface{}
}

type staticTableEntry struct {
	ID       int
	IconPath string
}

type fieldMapping struct {
	DBField   string
	JSONField string
}

var classes = map[string]staticTableEntry{
	"life":    {ID: 0, IconPath: "images/misc/class/life.png"},
	"death":   {ID: 1, IconPath: "images/misc/class/death.png"},
	"nature":  {ID: 2, IconPath: "images/misc/class/nature.png"},
	"sorcery": {ID: 3, IconPath: "images/misc/class/sorcery.png"},
	"chaos":   {ID: 4, IconPath: "images/misc/class/chaos.png"},
}

var stats = map[string]staticTableEntry{
	"health":       {ID: 0, IconPath: "images/misc/stat/health.png"},
	"attack":       {ID: 1, IconPath: "images/misc/stat/attack.png"},
	"intelligence": {ID: 2, IconPath: "images/misc/stat/intelligence.png"},
	"defense":      {ID: 3, IconPath: "images/misc/stat/defense.png"},
	"speed":        {ID: 4, IconPath: "images/misc/stat/speed.png"},
}

var jsonSources = [9]jsonMeta{
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

var correlatedFieldNamesMetaMap = map[string][]fieldMapping{
	"artifacts": {
		{DBField: "id", JSONField: "id"},
		{DBField: "name", JSONField: "name"},
		{DBField: "icon", JSONField: "icons"},
		{DBField: "stat_type", JSONField: "stat"},
	},
	"classes": {
		{DBField: "id", JSONField: ""},
		{DBField: "name", JSONField: ""},
		{DBField: "icon", JSONField: ""},
	},
	"races": {
		{DBField: "id", JSONField: ""},
		{DBField: "name", JSONField: "name"},
		{DBField: "icon", JSONField: "icon"},
	},
	"specializations": {
		{DBField: "id", JSONField: "id"},
		{DBField: "name", JSONField: "name"},
		{DBField: "description", JSONField: "description"},
		{DBField: "icon", JSONField: "icon"},
	},
	"materials": {
		{DBField: "id", JSONField: "id"},
		{DBField: "name", JSONField: "name"},
		{DBField: "description", JSONField: "description"},
		{DBField: "icon", JSONField: "icon"},
		{DBField: "type", JSONField: "slot"},
	},
	"traits": {
		{DBField: "id", JSONField: "id"},
		{DBField: "name", JSONField: "name"},
		{DBField: "description", JSONField: "description"},
		{DBField: "material_id", JSONField: "item"},
	},
	// creatures and material_stats use explicit mappers in mappers.go - not in this generic map
	"perks": {
		{DBField: "id", JSONField: "id"},
		{DBField: "name", JSONField: "name"},
		{DBField: "description", JSONField: "description"},
		{DBField: "icon", JSONField: "icon"},
		{DBField: "specialization_id", JSONField: "specialization"},
	},
	"spells": {
		{DBField: "id", JSONField: "id"},
		{DBField: "name", JSONField: "name"},
		{DBField: "description", JSONField: "description"},
		{DBField: "icon", JSONField: ""},
		{DBField: "charges", JSONField: "maxCharges"},
		{DBField: "class_id", JSONField: "class"},
	},
	"spellProperties": {
		{DBField: "id", JSONField: "id"},
		{DBField: "name", JSONField: "longDescription"},
		{DBField: "material_id", JSONField: "materialId"},
	},
	"stats": {
		{DBField: "id", JSONField: ""},
		{DBField: "stat_type", JSONField: ""},
	},
}

// TableRow represents a single row with column values
type TableRow map[string]interface{}

// Table represents all rows in a table, indexed by row ID
type Table map[int]TableRow

// TransformedData holds all transformed tables ready for DB insertion
type TransformedData struct {
	tables map[string]Table
}

// NewTransformedData creates an initialized TransformedData instance
func NewTransformedData() *TransformedData {
	return &TransformedData{
		tables: make(map[string]Table),
	}
}

// GetTable retrieves a table by name
func (td *TransformedData) GetTable(name string) (Table, bool) {
	table, exists := td.tables[name]
	return table, exists
}

// SetTable stores or replaces a table
func (td *TransformedData) SetTable(name string, table Table) {
	td.tables[name] = table
}

// EnsureTable ensures a table exists, creating it if necessary
func (td *TransformedData) EnsureTable(name string) Table {
	if td.tables[name] == nil {
		td.tables[name] = make(Table)
	}
	return td.tables[name]
}

// SetRow sets a row in a table, creating the table if needed
func (td *TransformedData) SetRow(tableName string, rowID int, row TableRow) {
	table := td.EnsureTable(tableName)
	table[rowID] = row
}

// GetRow retrieves a specific row from a table
func (td *TransformedData) GetRow(tableName string, rowID int) (TableRow, bool) {
	table, exists := td.tables[tableName]
	if !exists {
		return nil, false
	}
	row, exists := table[rowID]
	return row, exists
}

// Tables returns all tables (for compatibility with existing code)
func (td *TransformedData) Tables() map[string]Table {
	return td.tables
}

var correlatedTables *TransformedData

const gameDataRootPath = "/app/gameData"
