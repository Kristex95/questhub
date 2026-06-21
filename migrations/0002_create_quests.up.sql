-- Migration 0002: create quests table
CREATE TABLE IF NOT EXISTS quests (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description TEXT NOT NULL,
    difficulty INTEGER NOT NULL DEFAULT 1,
    xp_reward   INTEGER NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
