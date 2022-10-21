-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    login VARCHAR (255) NOT NULL UNIQUE,
    password VARCHAR (255) NOT NULL,
    balance INTEGER DEFAULT 0,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL
);
---- create above / drop below ----

DROP TABLE users;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
