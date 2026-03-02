// Package dbseeder - Explicit Mappers
//
// This file contains explicit type-safe mappers for tables with complex transformation logic.
// These mappers provide:
// - Compile-time type safety with typed input/output structs
// - Clearer error handling and validation
// - Easier debugging and testing
// - Better IDE support (autocomplete, refactoring)
//
// Simple tables continue to use config-driven generic mapping (see config.go).
// Use explicit mappers when:
// - Complex field transformations are needed
// - Multiple lookups or cross-table dependencies exist
// - Type conversions are non-trivial (e.g., arrays to specific indices)
// - Error handling needs to be more granular
package dbseeder

import (
	"log/slog"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// MapCreature transforms a creature from JSON format to DB format using SQLC-generated types
func MapCreature(raw map[string]interface{}, lookupTables *TransformedData) (repo.Creature, bool) {
	creature := repo.Creature{}

	// Extract ID
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

	// Extract name
	name, ok := raw["name"].(string)
	if !ok || name == "" {
		slog.Warn("Creature has invalid name", "id", id, "name", raw["name"])
		return creature, false
	}
	creature.Name = name

	// Extract and convert battle sprite to image bytes (optional - creature can exist without image)
	battleSprite, ok := raw["battleSprite"].(string)
	if ok && battleSprite != "" {
		imageBytes, ok := IconToBytes(battleSprite)
		if ok {
			creature.Image = imageBytes
		} else {
			slog.Debug("Failed to read creature image, continuing without image", "id", id, "battleSprite", battleSprite)
			creature.Image = nil
		}
	} else {
		slog.Debug("Creature has no battleSprite field, continuing without image", "id", id)
		creature.Image = nil
	}

	// Extract trait_id
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
	creature.TraitID = pgtype.Int4{Int32: int32(traitID), Valid: true}

	// Map class string to class_id
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

	// Map race string to race_id by looking up in races table
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

// ProcessCreaturesExplicit uses explicit mapping instead of generic field mapping
func ProcessCreaturesExplicit(creaturesJSON []map[string]interface{}, result *TransformedData) {
	creatureTable := result.EnsureTable("creatures")

	successCount := 0
	for _, rawCreature := range creaturesJSON {
		creature, ok := MapCreature(rawCreature, result)
		if !ok {
			continue
		}

		// Store in result map format expected by the pipeline
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

// MapMaterialToStats transforms material stats from JSON format to DB format using SQLC-generated types
// Returns a repo.MaterialStat entry with the first stat and optionally second stat
func MapMaterialToStats(raw map[string]interface{}, lookupTables *TransformedData) (repo.MaterialStat, bool) {
	entry := repo.MaterialStat{}

	// Extract material ID
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

	// Extract stats array
	statsRaw, exists := raw["stats"]
	if !exists {
		// Materials without stats are valid (e.g., trait materials)
		return entry, false
	}

	statsArray, ok := statsRaw.([]interface{})
	if !ok {
		slog.Warn("Material has invalid stats field type", "materialID", materialID, "stats", statsRaw)
		return entry, false
	}

	if len(statsArray) == 0 {
		// No stats to process
		return entry, false
	}

	// Extract stat names
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

	// Look up stats table
	statsTable, exists := lookupTables.GetTable("stats")
	if !exists {
		slog.Warn("Stats table not found", "materialID", materialID)
		return entry, false
	}

	// Resolve first stat name to ID
	firstStatID, found := findStatIDByName(statNames[0], statsTable)
	if !found {
		slog.Warn("Stat name not found", "materialID", materialID, "statName", statNames[0])
		return entry, false
	}
	entry.StatID = int32(firstStatID)

	// Resolve second stat name to ID if present
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
		if name, ok := fields["stat_type"].(string); ok && name == statName {
			return id, true
		}
	}
	return 0, false
}

// ProcessMaterialStatsExplicit uses explicit mapping instead of generic field mapping
func ProcessMaterialStatsExplicit(materialsJSON []map[string]interface{}, result *TransformedData) {
	materialStatsTable := result.EnsureTable("material_stats")

	idCounter := 0
	successCount := 0
	for _, rawMaterial := range materialsJSON {
		materialStat, ok := MapMaterialToStats(rawMaterial, result)
		if !ok {
			continue
		}

		// Store in result map format expected by the pipeline
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
