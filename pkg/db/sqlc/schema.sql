CREATE TABLE IF NOT EXISTS website (
        website_id INTEGER PRIMARY KEY AUTOINCREMENT,
        url        TEXT NOT NULL UNIQUE,
        created    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tag (
        tag_id     INTEGER PRIMARY KEY AUTOINCREMENT,
        website_id INT NOT NULL REFERENCES website (website_id),
        tag_type   TEXT,
        name       TEXT,
        property   TEXT,
        value      TEXT,
        created    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);