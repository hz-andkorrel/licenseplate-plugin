-- Migration 004: Create outbox table for reliable event delivery

CREATE TABLE IF NOT EXISTS outbox_events (
    id SERIAL PRIMARY KEY,
    channel VARCHAR(100) NOT NULL,
    payload TEXT NOT NULL,
    attempts INT NOT NULL DEFAULT 0,
    last_error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMP
);

CREATE INDEX idx_outbox_events_sent_at ON outbox_events(sent_at NULLS FIRST, created_at);

COMMENT ON TABLE outbox_events IS 'Outbox table for reliable event publishing to Redis';
