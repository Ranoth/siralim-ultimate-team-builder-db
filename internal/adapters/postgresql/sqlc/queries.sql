-- name: GetCreatures :many
SELECT c.id,
    c.name,
    c.icon,
    t.name as trait,
    cl.name as class,
    r.name as race
FROM creatures c
    LEFT JOIN traits t ON c.trait_id = t.id
    LEFT JOIN classes cl ON c.class_id = cl.id
    LEFT JOIN races r ON c.race_id = r.id;
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
SELECT
    c.id,
    c.name,
    c.icon AS icon,
    t.name AS trait,
    cl.name AS class,
    r.name AS race,
    jsonb_object_agg(s.type, csg.growth_rate) AS stats
FROM creatures c
JOIN classes cl
    ON c.class_id = cl.id
JOIN races r
    ON c.race_id = r.id
LEFT JOIN traits t
    ON c.trait_id = t.id
JOIN creature_stat_growth csg
    ON c.id = csg.creature_id
JOIN stats s
    ON csg.stat_id = s.id
WHERE c.id = $1
GROUP BY
    c.id,
    c.name,
    c.icon,
    t.name,
    cl.name,
    r.name;
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
-- name: GetMaterialStats :many
SELECT id,
    material_id,
    stat_id
FROM material_stats
WHERE material_id = $1;
-- name: GetStatGrowthsByCreatureId :many
SELECT csg.id,
    c.name as creature_name,
    s.type as stat_type,
    csg.growth_rate
FROM creature_stat_growth csg
    JOIN stats s ON csg.stat_id = s.id
    JOIN creatures c ON csg.creature_id = c.id
WHERE csg.creature_id = $1;
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
-- name: GetRelics :many
SELECT r.id,
    r.name,
    r.icon,
    r.bonuses,
    s.id as stat_id,
    s.type as stat_type
FROM relics r
    LEFT JOIN stats s ON r.stat_id = s.id;
-- name: GetRelic :one
SELECT r.id,
    r.name,
    r.icon,
    r.bonuses,
    s.id as stat_id,
    s.type as stat_type
FROM relics r
    LEFT JOIN stats s ON r.stat_id = s.id
WHERE r.id = $1;
-- name: GetRelicsByName :many
SELECT r.id,
    r.name,
    r.icon,
    r.bonuses,
    s.id as stat_id,
    s.type as stat_type
FROM relics r
    LEFT JOIN stats s ON r.stat_id = s.id
WHERE r.name ILIKE '%' || $1 || '%';
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