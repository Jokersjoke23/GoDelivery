CREATE TABLE IF NOT EXISTS exports (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type        VARCHAR(20) NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    file_path   VARCHAR(255),
    filters     TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);