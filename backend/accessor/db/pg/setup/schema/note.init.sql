DO $$
BEGIN
    CREATE TYPE note_status AS ENUM ('Active', 'Archived');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS note (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    status note_status NOT NULL DEFAULT 'Active',
    created_ts TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_updated_ts TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_note_status_id ON note(status, id);
