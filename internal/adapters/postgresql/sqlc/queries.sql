-- name: GetCreatures :many
SELECT *
FROM creatures_view;
-- name: GetTraits :many
SELECT *
FROM traits;
-- name: GetClasses :many
SELECT *
FROM classes_view;
-- name: GetRaces :many
SELECT *
FROM races_view;
-- name: GetSpecializations :many
SELECT *
FROM specializations;
-- name: GetPerks :many
SELECT *
FROM perks;
-- name: GetSpells :many
SELECT *
FROM spells;
-- name: GetMaterials :many
SELECT *
FROM materials_view;
-- name: GetSpellProperties :many
SELECT *
FROM spell_properties;
-- name: GetArtifacts :many
SELECT *
FROM artifacts;
-- name: GetStats :many
SELECT *
FROM stats;
-- name: GetCreature :one
SELECT *
FROM creatures_view
WHERE id = $1;
-- name: GetCreatureIconById :one
SELECT icon
FROM creatures
WHERE id = $1;
-- name: GetTrait :one
SELECT *
FROM traits
WHERE id = $1;
-- name: GetClass :one
SELECT *
FROM classes_view
WHERE id = $1;
-- name: GetRace :one
SELECT *
FROM races_view
WHERE id = $1;
-- name: GetSpecialization :one
SELECT *
FROM specializations
WHERE id = $1;
-- name: GetPerk :one
SELECT *
FROM perks
WHERE id = $1;
-- name: GetSpell :one
SELECT *
FROM spells
WHERE id = $1;
-- name: GetMaterial :one
SELECT *
FROM materials_view
WHERE id = $1;
-- name: GetMaterialIconById :one
SELECT icon
FROM materials
WHERE id = $1;
-- name: MaterialExists :one
SELECT EXISTS(
        SELECT 1
        FROM materials
        WHERE id = $1
    ) as exists;
-- name: TraitExists :one
SELECT EXISTS(
        SELECT 1
        FROM traits
        WHERE id = $1
    ) as exists;
-- name: GetSpellProperty :one
SELECT *
FROM spell_properties
WHERE id = $1;
-- name: GetArtifact :one
SELECT *
FROM artifacts
WHERE id = $1;
-- name: GetStat :one
SELECT *
FROM stats
WHERE id = $1;
-- name: GetStatsCount :one
SELECT COUNT(*)
FROM stats;
-- name: GetTraitsByCreatureName :many
SELECT t.*
FROM traits t
    JOIN creatures c ON t.id = c.trait_id
WHERE c.name ILIKE '%' || $1 || '%';
-- name: GetTraitsByName :many
SELECT *
FROM traits
WHERE name ILIKE '%' || $1 || '%';
-- name: GetClassesByName :many
SELECT *
FROM classes_view
WHERE name ILIKE '%' || $1 || '%';
-- name: GetRacesByName :many
SELECT *
FROM races_view
WHERE name ILIKE '%' || $1 || '%';
-- name: GetSpecializationsByName :many
SELECT *
FROM specializations
WHERE name ILIKE '%' || $1 || '%';
-- name: GetPerksByName :many
SELECT *
FROM perks
WHERE name ILIKE '%' || $1 || '%';
-- name: GetSpellsByName :many
SELECT *
FROM spells
WHERE name ILIKE '%' || $1 || '%';
-- name: GetMaterialsByName :many
SELECT *
FROM materials_view
WHERE name ILIKE '%' || $1 || '%';
-- name: GetSpellPropertiesByName :many
SELECT *
FROM spell_properties
WHERE name ILIKE '%' || $1 || '%';
-- name: GetArtifactsByName :many
SELECT *
FROM artifacts
WHERE name ILIKE '%' || $1 || '%';
-- name: GetStatsByType :many
SELECT *
FROM stats
WHERE type::text ILIKE '%' || $1 || '%';
-- name: GetCreaturesByTraitName :many
SELECT *
FROM creatures_view
WHERE trait ILIKE '%' || $1 || '%';
-- name: GetCreaturesByClassName :many
SELECT *
FROM creatures_view
WHERE class ILIKE '%' || $1 || '%';
-- name: GetCreaturesByRaceName :many
SELECT *
FROM creatures_view
WHERE race ILIKE '%' || $1 || $1 || '%';
-- name: GetCreaturesByName :many
SELECT *
FROM creatures_view
WHERE name ILIKE '%' || $1 || '%';
-- name: GetRacesByTraitName :many
SELECT r.*
FROM races_view r
    JOIN creatures c ON r.id = c.race_id
    JOIN traits t ON c.trait_id = t.id
WHERE t.name ILIKE '%' || $1 || '%';
-- name: GetRacesByClassName :many
SELECT r.*
FROM races_view r
    JOIN creatures c ON r.id = c.race
    JOIN classes cl ON c.class_id = cl.id
WHERE cl.name ILIKE '%' || $1 || '%';
-- name: GetRacesByCreatureName :many
SELECT r.*
FROM races_view r
    JOIN creatures c ON r.id = c.race_id
WHERE c.name ILIKE '%' || $1 || '%';
-- name: GetRelics :many
SELECT *
FROM relics_view;
-- name: GetRelic :one
SELECT *
FROM relics_view
WHERE id = $1;
-- name: GetRelicsByName :many
SELECT *
FROM relics_view
WHERE name ILIKE '%' || $1 || '%';
-- name: GetRelicIconById :one
SELECT icon
FROM relics
WHERE id = $1;
-- name: GetRaceIconById :one
SELECT icon
FROM races
WHERE id = $1;
-- name: GetClassIconById :one
SELECT icon
FROM classes
WHERE id = $1;
-- Batch insert queries using COPY protocol for efficient seeding
-- name: BatchInsertClasses :copyfrom
INSERT INTO classes (id, name, icon)
VALUES ($1, $2, $3);
-- name: BatchInsertStats :copyfrom
INSERT INTO stats (id, type, icon)
VALUES ($1, $2, $3);
-- name: BatchInsertRaces :copyfrom
INSERT INTO races (id, name, icon)
VALUES ($1, $2, $3);
-- name: BatchInsertSpecializations :copyfrom
INSERT INTO specializations (id, name, description, icon)
VALUES ($1, $2, $3, $4);
-- name: BatchInsertMaterials :copyfrom
INSERT INTO materials (id, name, icon, type)
VALUES ($1, $2, $3, $4);
-- name: BatchInsertTraits :copyfrom
INSERT INTO traits (id, name, description, material_id)
VALUES ($1, $2, $3, $4);
-- name: BatchInsertPerks :copyfrom
INSERT INTO perks (id, name, description, icon, specialization_id)
VALUES ($1, $2, $3, $4, $5);
-- name: BatchInsertSpells :copyfrom
INSERT INTO spells (id, name, description, charges, class_id)
VALUES ($1, $2, $3, $4, $5);
-- name: BatchInsertSpellProperties :copyfrom
INSERT INTO spell_properties (id, name, material_id)
VALUES ($1, $2, $3);
-- name: BatchInsertArtifacts :copyfrom
INSERT INTO artifacts (id, name, icon, stat_id)
VALUES ($1, $2, $3, $4);
-- name: BatchInsertCreatures :copyfrom
INSERT INTO creatures (id, name, icon, trait_id, class_id, race_id)
VALUES ($1, $2, $3, $4, $5, $6);
-- name: BatchInsertRelics :copyfrom
INSERT INTO relics (id, name, icon, bonuses, stat_id)
VALUES ($1, $2, $3, $4, $5);
-- name: BatchInsertMaterialStats :copyfrom
INSERT INTO material_stats (id, material_id, stat_id, stat_id2)
VALUES ($1, $2, $3, $4);
-- name: BatchInsertCreatureStatGrowth :copyfrom
INSERT INTO creature_stat_growth (id, creature_id, stat_id, growth_rate)
VALUES ($1, $2, $3, $4);