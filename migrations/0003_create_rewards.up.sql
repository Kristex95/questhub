-- Migration 0003: create rewards table
CREATE TABLE IF NOT EXISTS rewards (
    id SERIAL PRIMARY KEY,
    quest_id INTEGER NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    reward_type VARCHAR NOT NULL,
    xp_amount INTEGER NOT NULL DEFAULT 0,
    item_name VARCHAR NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
