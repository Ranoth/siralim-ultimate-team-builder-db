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
    id SERIAL PRIMARY KEY,
    type stat_type NOT NULL
);
CREATE TABLE IF NOT EXISTS materials (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    icon BYTEA NOT NULL,
    type material_type DEFAULT 'trait'
);
CREATE TABLE IF NOT EXISTS material_stats (
    id SERIAL PRIMARY KEY,
    material_id INTEGER NOT NULL,
    stat_id INTEGER NOT NULL,
    stat_id2 INTEGER,
    FOREIGN KEY (material_id) REFERENCES materials(id),
    FOREIGN KEY (stat_id) REFERENCES stats(id),
    FOREIGN KEY (stat_id2) REFERENCES stats(id)
);
CREATE TABLE IF NOT EXISTS artifacts (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon BYTEA NOT NULL,
    type stat_type NOT NULL
);
CREATE TABLE IF NOT EXISTS spell_properties (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    material_id INTEGER NOT NULL,
    FOREIGN KEY (material_id) REFERENCES materials(id)
);
CREATE TABLE IF NOT EXISTS traits (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    material_id INTEGER NOT NULL,
    FOREIGN KEY (material_id) REFERENCES materials(id)
);
CREATE TABLE IF NOT EXISTS classes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    icon BYTEA NOT NULL
);
CREATE TABLE IF NOT EXISTS races (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    icon BYTEA NOT NULL
);
CREATE TABLE IF NOT EXISTS creatures (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    image BYTEA NOT NULL,
    trait_id INTEGER NOT NULL,
    class_id INTEGER NOT NULL,
    race_id INTEGER NOT NULL,
    FOREIGN KEY (trait_id) REFERENCES traits(id),
    FOREIGN KEY (class_id) REFERENCES classes(id),
    FOREIGN KEY (race_id) REFERENCES races(id)
);
CREATE TABLE IF NOT EXISTS specializations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon BYTEA NOT NULL
);
CREATE TABLE IF NOT EXISTS perks (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon BYTEA NOT NULL,
    specialization_id INTEGER NOT NULL,
    FOREIGN KEY (specialization_id) REFERENCES specializations(id)
);
CREATE TABLE IF NOT EXISTS spells (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon BYTEA NOT NULL,
    charges INTEGER NOT NULL,
    class_id INTEGER NOT NULL,
    FOREIGN KEY (class_id) REFERENCES classes(id)
);
-- +goose Down
DROP TABLE IF EXISTS materials CASCADE;
DROP TABLE IF EXISTS artifacts CASCADE;
DROP TABLE IF EXISTS spell_properties CASCADE;
DROP TABLE IF EXISTS stats CASCADE;
DROP TABLE IF EXISTS material_stats CASCADE;
DROP TABLE IF EXISTS spells CASCADE;
DROP TABLE IF EXISTS perks CASCADE;
DROP TABLE IF EXISTS specializations CASCADE;
DROP TABLE IF EXISTS creatures CASCADE;
DROP TABLE IF EXISTS races CASCADE;
DROP TABLE IF EXISTS classes CASCADE;
DROP TABLE IF EXISTS traits CASCADE;
DROP TYPE IF EXISTS material_type CASCADE;
DROP TYPE IF EXISTS stat_type CASCADE;