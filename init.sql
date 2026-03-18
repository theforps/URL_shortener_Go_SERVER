CREATE TABLE IF NOT EXISTS url_table (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    uniq_code TEXT,
    url_base TEXT,
    views INTEGER,
    finally_date TEXT
);

CREATE INDEX idx_url_table_uniq_code
ON url_table (uniq_code);