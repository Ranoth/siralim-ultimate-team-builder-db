package dbseeder

import (
	"log/slog"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func MapCreature(raw map[string]interface{}, lookupTables *TransformedData) (repo.Creature, bool) {
	creature := repo.Creature{}

	idRaw, exists := raw["id"]
	if !exists {
		slog.Warn("Creature missing id field", "raw", raw)
		return creature, false
	}
	id, ok := GetIDValue(idRaw)
	if !ok {
		slog.Warn("Creature has invalid id", "idRaw", idRaw)
		return creature, false
	}
	creature.ID = int32(id)

	name, ok := raw["name"].(string)
	if !ok || name == "" {
		slog.Warn("Creature has invalid name", "id", id, "name", raw["name"])
		return creature, false
	}
	creature.Name = name

	battleSprite, ok := raw["battleSprite"].(string)
	if !ok || battleSprite == "" {
		slog.Warn("Creature has invalid battleSprite", "id", id, "battleSprite", raw["battleSprite"])
		return creature, false
	}
	imageBytes, ok := IconToBytes(battleSprite)
	if !ok {
		slog.Warn("Failed to read creature image", "id", id, "battleSprite", battleSprite)
		return creature, false
	}
	creature.Image = imageBytes

	traitRaw, exists := raw["trait"]
	if !exists {
		slog.Warn("Creature missing trait field", "id", id)
		return creature, false
	}
	traitID, ok := GetIDValue(traitRaw)
	if !ok {
		slog.Warn("Creature has invalid trait", "id", id, "trait", traitRaw)
		return creature, false
	}
	creature.TraitID = int32(traitID)

	className, ok := raw["class"].(string)
	if !ok || className == "" {
		slog.Warn("Creature has invalid class", "id", id, "class", raw["class"])
		return creature, false
	}
	classEntry, exists := classes[className]
	if !exists {
		slog.Warn("Unknown class name", "id", id, "class", className)
		return creature, false
	}
	creature.ClassID = int32(classEntry.ID)

	raceName, ok := raw["race"].(string)
	if !ok || raceName == "" {
		slog.Warn("Creature has invalid race", "id", id, "race", raw["race"])
		return creature, false
	}
	racesTable, exists := lookupTables.GetTable("races")
	if !exists {
		slog.Warn("Races table not found", "id", id)
		return creature, false
	}
	raceID, found := findRaceIDByName(raceName, racesTable)
	if !found {
		slog.Warn("Unknown race name", "id", id, "race", raceName)
		return creature, false
	}
	creature.RaceID = int32(raceID)

	return creature, true
}

func findRaceIDByName(raceName string, racesTable Table) (int, bool) {
	for id, fields := range racesTable {
		if name, ok := fields["name"].(string); ok && name == raceName {
			return id, true
		}
	}
	return 0, false
}

func ProcessCreaturesExplicit(creaturesJSON []map[string]interface{}, result *TransformedData) {
	creatureTable := result.EnsureTable("creatures")

	successCount := 0
	for _, rawCreature := range creaturesJSON {
		creature, ok := MapCreature(rawCreature, result)
		if !ok {
			continue
		}

		creatureTable[int(creature.ID)] = TableRow{
			"id":       creature.ID,
			"name":     creature.Name,
			"image":    creature.Image,
			"trait_id": creature.TraitID,
			"class_id": creature.ClassID,
			"race_id":  creature.RaceID,
		}
		successCount++
	}

	slog.Info("Processed creatures with explicit mapper", "total", len(creaturesJSON), "success", successCount, "failed", len(creaturesJSON)-successCount)
}

func MapMaterialToStats(raw map[string]interface{}, lookupTables *TransformedData) (repo.MaterialStat, bool) {
	entry := repo.MaterialStat{}

	idRaw, exists := raw["id"]
	if !exists {
		slog.Warn("Material missing id field", "raw", raw)
		return entry, false
	}
	materialID, ok := GetIDValue(idRaw)
	if !ok {
		slog.Warn("Material has invalid id", "idRaw", idRaw)
		return entry, false
	}
	entry.MaterialID = int32(materialID)

	statsRaw, exists := raw["stats"]
	if !exists {
		return entry, false
	}

	statsArray, ok := statsRaw.([]interface{})
	if !ok {
		slog.Warn("Material has invalid stats field type", "materialID", materialID, "stats", statsRaw)
		return entry, false
	}

	if len(statsArray) == 0 {
		return entry, false
	}

	statNames := make([]string, 0, len(statsArray))
	for _, statRaw := range statsArray {
		statName, ok := statRaw.(string)
		if !ok {
			slog.Warn("Material has invalid stat name type", "materialID", materialID, "statName", statRaw)
			continue
		}
		statNames = append(statNames, statName)
	}

	if len(statNames) == 0 {
		return entry, false
	}

	statsTable, exists := lookupTables.GetTable("stats")
	if !exists {
		slog.Warn("Stats table not found", "materialID", materialID)
		return entry, false
	}

	firstStatID, found := findStatIDByName(statNames[0], statsTable)
	if !found {
		slog.Warn("Stat name not found", "materialID", materialID, "statName", statNames[0])
		return entry, false
	}
	entry.StatID = int32(firstStatID)

	if len(statNames) >= 2 {
		secondStatID, found := findStatIDByName(statNames[1], statsTable)
		if found {
			entry.StatId2 = pgtype.Int4{Int32: int32(secondStatID), Valid: true}
		}
	}

	return entry, true
}

func findStatIDByName(statName string, statsTable Table) (int, bool) {
	for id, fields := range statsTable {
		if name, ok := fields["name"].(string); ok && name == statName {
			return id, true
		}
	}
	return 0, false
}

func ProcessMaterialStatsExplicit(materialsJSON []map[string]interface{}, result *TransformedData) {
	materialStatsTable := result.EnsureTable("material_stats")

	idCounter := 0
	successCount := 0
	for _, rawMaterial := range materialsJSON {
		materialStat, ok := MapMaterialToStats(rawMaterial, result)
		if !ok {
			continue
		}

		materialStatsTable[idCounter] = TableRow{
			"id":          idCounter,
			"material_id": materialStat.MaterialID,
			"stat_id":     materialStat.StatID,
			"stat_id2":    materialStat.StatId2,
		}
		idCounter++
		successCount++
	}

	slog.Info("Processed material_stats with explicit mapper", "total", len(materialsJSON), "success", successCount, "failed", len(materialsJSON)-successCount)
}
