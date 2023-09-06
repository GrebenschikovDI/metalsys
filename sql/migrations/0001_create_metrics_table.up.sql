CREATE TABLE IF NOT EXISTS metrics
(
    id    TEXT PRIMARY KEY,
    mtype TEXT,
    delta BIGINT,
    value DOUBLE PRECISION
);