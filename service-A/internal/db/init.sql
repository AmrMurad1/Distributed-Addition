-- Create outbox table for implementing outbox pattern
CREATE TABLE IF NOT EXISTS outbox (
     id SERIAL PRIMARY KEY,
     event_type VARCHAR(100) NOT NULL,
     payload JSONB NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     processed BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_outbox_processed_created
    ON outbox (processed, created_at) WHERE processed = FALSE;