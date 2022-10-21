-- Write your migrate up statements here
CREATE TYPE status AS ENUM('processed', 'processing', 'invalid', 'new');
CREATE TABLE IF NOT EXISTS orders(
   id BIGSERIAL PRIMARY KEY,
   number VARCHAR(255) NOT NULL UNIQUE,
   status status NOT NULL,
   amount BIGINT NOT NULL,
   user_id INTEGER NOT NULL,
   created_at TIMESTAMP NULL,
   updated_at TIMESTAMP NULL,
   CONSTRAINT fk_orders_users FOREIGN KEY (user_id) REFERENCES users (id)
);
---- create above / drop below ----

DROP TABLE orders;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
