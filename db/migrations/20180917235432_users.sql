-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS users (
       id SERIAL PRIMARY KEY,
       username VARCHAR(255) NOT NULL DEFAULT '',
       firstname VARCHAR(255) NOT NULL DEFAULT '',
       lastname VARCHAR(255) NOT NULL DEFAULT '',
       password VARCHAR(255) NOT NULL DEFAULT '',
       email VARCHAR(255) NOT NULL DEFAULT '',
       role VARCHAR(255) NOT NULL DEFAULT 'customer',
       uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
       CONSTRAINT username UNIQUE (username),
       CONSTRAINT email UNIQUE (email)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "users";

