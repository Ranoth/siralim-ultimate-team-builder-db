-- name: GetCreatures :many
SELECT *
FROM creatures;
-- name: GetTraits :many
SELECT *
FROM traits;
-- name: GetClasses :many
SELECT *
FROM classes;
-- name: GetRaces :many
SELECT *
FROM races;
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
SELECT m.id,
    m.name,
    m.icon,
    m.type,
    ms.id as stat_id,
    ms.stat_id
FROM materials m
    LEFT JOIN material_stats ms ON m.id = ms.material_id;
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
FROM creatures
WHERE id = $1;
-- name: GetTrait :one
SELECT *
FROM traits
WHERE id = $1;
-- name: GetClass :one
SELECT *
FROM classes
WHERE id = $1;
-- name: GetRace :one
SELECT *
FROM races
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
SELECT m.id,
    m.name,
    m.icon,
    m.type,
    ms.id as stat_id,
    ms.stat_id
FROM materials m
    LEFT JOIN material_stats ms ON m.id = ms.material_id
WHERE m.id = $1;
-- name: MaterialExists :one
SELECT EXISTS(SELECT 1 FROM materials WHERE id = $1) as exists;
-- name: TraitExists :one
SELECT EXISTS(SELECT 1 FROM traits WHERE id = $1) as exists;
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
-- name: CreateCreature :one
INSERT INTO creatures (id, name, image, trait_id, class_id, race_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
-- name: CreateTrait :one
INSERT INTO traits (id, name, description, material_id)
VALUES ($1, $2, $3, $4)
RETURNING id;
-- name: CreateClass :one
INSERT INTO classes (id, name, icon)
VALUES ($1, $2, $3)
RETURNING id;
-- name: CreateRace :one
INSERT INTO races (id, name, icon)
VALUES ($1, $2, $3)
RETURNING id;
-- name: CreateSpecialization :one
INSERT INTO specializations (id, name, description)
VALUES ($1, $2, $3)
RETURNING id;
-- name: CreatePerk :one
INSERT INTO perks (id, name, description, icon, specialization_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
-- name: CreateSpell :one
INSERT INTO spells (id, name, description, icon, charges, class_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
-- name: CreateMaterial :one
INSERT INTO materials (id, name, icon, type)
VALUES ($1, $2, $3, $4)
RETURNING id,
    name,
    icon,
    type;
-- name: CreateSpellProperty :one
INSERT INTO spell_properties (id, name, material_id)
VALUES ($1, $2, $3)
RETURNING id;
-- name: CreateArtifact :one
INSERT INTO artifacts (id, name, icon, type)
VALUES ($1, $2, $3, $4)
RETURNING id;
-- name: CreateStat :one
INSERT INTO stats (id, type)
VALUES ($1, $2)
RETURNING id;
-- name: GetStatsCount :one
SELECT COUNT(*) FROM stats;
-- name: DeleteCreature :exec
DELETE FROM creatures
WHERE id = $1;
-- name: DeleteTrait :exec
DELETE FROM traits
WHERE id = $1;
-- name: DeleteClass :exec
DELETE FROM classes
WHERE id = $1;
-- name: DeleteRace :exec
DELETE FROM races
WHERE id = $1;
-- name: DeleteSpecialization :exec
DELETE FROM specializations
WHERE id = $1;
-- name: DeletePerk :exec
DELETE FROM perks
WHERE id = $1;
-- name: DeleteSpell :exec
DELETE FROM spells
WHERE id = $1;
-- name: DeleteMaterial :exec
DELETE FROM materials
WHERE id = $1;
-- name: GetMaterialStats :many
SELECT id,
    material_id,
    stat_id
FROM material_stats
WHERE material_id = $1;
-- name: CreateMaterialStat :one
INSERT INTO material_stats (material_id, stat_id, stat_id2, id)
VALUES ($1, $2, $3, $4)
RETURNING id;
-- name: UpdateMaterialStat :exec
UPDATE material_stats
SET id = $3
WHERE material_id = $1
    AND stat_id = $2;
-- name: DeleteMaterialStat :exec
DELETE FROM material_stats
WHERE id = $1;
-- name: DeleteSpellProperty :exec
DELETE FROM spell_properties
WHERE id = $1;
-- name: DeleteArtifact :exec
DELETE FROM artifacts
WHERE id = $1;
-- name: DeleteStat :exec
DELETE FROM stats
WHERE id = $1;
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
FROM classes
WHERE name ILIKE '%' || $1 || '%';
-- name: GetRacesByName :many
SELECT *
FROM races
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
SELECT m.id,
    m.name,
    m.icon,
    m.type,
    ms.id as stat_id,
    ms.stat_id
FROM materials m
    LEFT JOIN material_stats ms ON m.id = ms.material_id
WHERE m.name ILIKE '%' || $1 || '%';
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
SELECT c.*
FROM creatures c
    JOIN traits t ON c.trait_id = t.id
WHERE t.name ILIKE '%' || $1 || '%';
-- name: GetCreaturesByClassName :many
SELECT c.*
FROM creatures c
    JOIN classes cl ON c.class_id = cl.id
WHERE cl.name ILIKE '%' || $1 || '%';
-- name: GetCreaturesByRaceName :many
SELECT c.*
FROM creatures c
    JOIN races r ON c.race_id = r.id
WHERE r.name ILIKE '%' || $1 || '%';
-- name: GetCreaturesByName :many
SELECT *
FROM creatures
WHERE name ILIKE '%' || $1 || '%';
-- name: GetRacesByTraitName :many
SELECT r.*
FROM races r
    JOIN creatures c ON r.id = c.race_id
    JOIN traits t ON c.trait_id = t.id
WHERE t.name ILIKE '%' || $1 || '%';
-- name: GetRacesByClassName :many
SELECT r.*
FROM races r
    JOIN creatures c ON r.id = c.race
    JOIN classes cl ON c.class_id = cl.id
WHERE cl.name ILIKE '%' || $1 || '%';
-- name: GetRacesByCreatureName :many
SELECT r.*
FROM races r
    JOIN creatures c ON r.id = c.race_id
WHERE c.name ILIKE '%' || $1 || '%';