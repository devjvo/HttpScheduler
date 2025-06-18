CREATE TABLE IF NOT EXISTS request(
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    http_method VARCHAR(7) NOT NULL,
    url TEXT NOT NULL,
    response_code INT
);
