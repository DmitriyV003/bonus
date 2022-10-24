CREATE TYPE transaction_type AS ENUM('debit', 'credit');
CREATE TYPE type AS ENUM('withdraw', 'accrual');

CREATE TABLE IF NOT EXISTS payments(
   id SERIAL PRIMARY KEY,
   user_id INTEGER NOT NULL,
   type type NOT NULL,
   transaction_type transaction_type NOT NULL,
   order_number VARCHAR(255) NOT NULL,
   amount INTEGER NOT NULL,
   created_at TIMESTAMP NULL,
   updated_at TIMESTAMP NULL,
   CONSTRAINT fk_payments_users FOREIGN KEY (user_id) REFERENCES users (id)
);
---- create above / drop below ----

DROP TABLE payments;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
