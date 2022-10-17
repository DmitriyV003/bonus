-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS orders(
   id SERIAL PRIMARY KEY,
   number INTEGER NOT NULL,
   status VARCHAR(255),
   amount INTEGER,
   user_id INTEGER NOT NULL,
   created_at TIMESTAMP NULL,
   updated_at TIMESTAMP NULL,
   CONSTRAINT fk_orders_users FOREIGN KEY (user_id) REFERENCES users (id)
);
---- create above / drop below ----

DROP TABLE orders;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
