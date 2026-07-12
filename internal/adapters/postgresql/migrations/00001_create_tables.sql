-- +goose Up
CREATE TYPE material_type AS ENUM ('stat', 'trick', 'trait');
CREATE TYPE stat_type AS ENUM (
    'health',
    'attack',
    'intelligence',
    'defense',
    'speed'
);
CREATE TABLE IF NOT EXISTS stats (
    id INTEGER PRIMARY KEY,
    type stat_type NOT NULL,
    icon BYTEA NOT NULL
);
CREATE TABLE IF NOT EXISTS materials (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    icon BYTEA NOT NULL,
    type material_type NOT NULL
);
CREATE TABLE IF NOT EXISTS material_stats (
    id INTEGER PRIMARY KEY,
    material_id INTEGER NOT NULL,
    stat_id INTEGER NOT NULL,
    stat_id2 INTEGER,
    FOREIGN KEY (material_id) REFERENCES materials(id),
    FOREIGN KEY (stat_id) REFERENCES stats(id),
    FOREIGN KEY (stat_id2) REFERENCES stats(id)
);
CREATE TABLE IF NOT EXISTS artifacts (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    icon BYTEA NOT NULL,
    stat_id INTEGER NOT NULL,
    FOREIGN KEY (stat_id) REFERENCES stats(id)
);
CREATE TABLE IF NOT EXISTS spell_properties (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    material_id INTEGER NOT NULL,
    FOREIGN KEY (material_id) REFERENCES materials(id)
);
CREATE TABLE IF NOT EXISTS traits (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    material_id INTEGER,
    FOREIGN KEY (material_id) REFERENCES materials(id)
);
CREATE TABLE IF NOT EXISTS classes (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    icon BYTEA NOT NULL
);
CREATE TABLE IF NOT EXISTS races (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    icon BYTEA
);
CREATE TABLE IF NOT EXISTS creatures (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    icon BYTEA,
    trait_id INTEGER,
    class_id INTEGER NOT NULL,
    race_id INTEGER NOT NULL,
    FOREIGN KEY (trait_id) REFERENCES traits(id),
    FOREIGN KEY (class_id) REFERENCES classes(id),
    FOREIGN KEY (race_id) REFERENCES races(id)
);
CREATE TABLE IF NOT EXISTS creature_stat_growth (
    id INTEGER PRIMARY KEY,
    creature_id INTEGER NOT NULL,
    stat_id INTEGER NOT NULL,
    growth_rate INTEGER NOT NULL,
    FOREIGN KEY (creature_id) REFERENCES creatures(id),
    FOREIGN KEY (stat_id) REFERENCES stats(id)
);
CREATE TABLE IF NOT EXISTS specializations (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon BYTEA
);
CREATE TABLE IF NOT EXISTS perks (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon BYTEA NOT NULL,
    specialization_id INTEGER NOT NULL,
    FOREIGN KEY (specialization_id) REFERENCES specializations(id)
);
CREATE TABLE IF NOT EXISTS spells (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    charges INTEGER NOT NULL,
    class_id INTEGER NOT NULL,
    FOREIGN KEY (class_id) REFERENCES classes(id)
);
CREATE TABLE IF NOT EXISTS relics (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    icon BYTEA,
    bonuses text [] NOT NULL,
    stat_id INTEGER NOT NULL,
    FOREIGN KEY (stat_id) REFERENCES stats(id)
);
CREATE VIEW creatures_view AS
SELECT c.id,
    c.name,
    '/icons/creatures/' || c.id::text AS icon,
    t.name AS trait,
    t.description AS trait_description,
    cl.name AS class,
    r.name AS race,
    jsonb_object_agg(s.type, csg.growth_rate) AS stats
FROM creatures c
    JOIN classes cl ON c.class_id = cl.id
    JOIN races r ON c.race_id = r.id
    LEFT JOIN traits t ON c.trait_id = t.id
    JOIN creature_stat_growth csg ON c.id = csg.creature_id
    JOIN stats s ON csg.stat_id = s.id
GROUP BY c.id,
    c.name,
    c.icon,
    t.name,
    t.description,
    cl.name,
    r.name;
CREATE VIEW materials_view AS
SELECT m.id,
    m.name,
    '/icons/materials/' || m.id::text AS icon,
    m.type,
    ms.id as stat_id1,
    ms.stat_id2
FROM materials m
    LEFT JOIN material_stats ms ON m.id = ms.material_id;
CREATE VIEW relics_view AS
SELECT r.id,
    r.name,
    '/icons/relics/' || r.id::text AS icon,
    r.bonuses,
    s.id as stat_id,
    s.type as stat_type
FROM relics r
    LEFT JOIN stats s ON r.stat_id = s.id;
CREATE VIEW races_view AS
SELECT id,
    name,
    '/icons/races/' || id::text AS icon
FROM races;
-- +goose Down
DROP TABLE IF EXISTS relics CASCADE;
DROP TABLE IF EXISTS materials CASCADE;
DROP TABLE IF EXISTS artifacts CASCADE;
DROP TABLE IF EXISTS spell_properties CASCADE;
DROP TABLE IF EXISTS stats CASCADE;
DROP TABLE IF EXISTS material_stats CASCADE;
DROP TABLE IF EXISTS spells CASCADE;
DROP TABLE IF EXISTS perks CASCADE;
DROP TABLE IF EXISTS specializations CASCADE;
DROP TABLE IF EXISTS creatures CASCADE;
DROP TABLE IF EXISTS creature_stat_growth CASCADE;
DROP TABLE IF EXISTS races CASCADE;
DROP TABLE IF EXISTS classes CASCADE;
DROP TABLE IF EXISTS traits CASCADE;
DROP TYPE IF EXISTS material_type CASCADE;
DROP TYPE IF EXISTS stat_type CASCADE;
DROP VIEW IF EXISTS creatures_view CASCADE;
DROP VIEW IF EXISTS materials_view CASCADE;
DROP VIEW IF EXISTS relics_view CASCADE;
DROP VIEW IF EXISTS races_view CASCADE;