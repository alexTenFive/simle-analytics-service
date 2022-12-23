CREATE TABLE timestamps (
    id SERIAL PRIMARY KEY,
    ts TIMESTAMP NOT NULL,
    value INTEGER DEFAULT(0) NOT NULL
);
CREATE INDEX IF NOT EXISTS timestamps_time_id_idx ON timestamps (ts);

