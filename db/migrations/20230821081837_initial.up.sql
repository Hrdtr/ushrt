CREATE TABLE links (
    id VARCHAR(255) PRIMARY KEY,
    original_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX unique_original_url ON links (original_url);
