-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS trains (
       id SERIAL PRIMARY KEY,
       brandname VARCHAR(255) NOT NULL DEFAULT '',
       brandlogo VARCHAR(255) NOT NULL DEFAULT '',
       brandfeatures VARCHAR(255) NOT NULL DEFAULT '',
       CONSTRAINT brandname UNIQUE (brandname)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "trains";

