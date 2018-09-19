
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS routes (
       id SERIAL PRIMARY KEY,
       routes_from varchar(255) NOT NULL DEFAULT '',
       brandname varchar(255) NOT NULL DEFAULT '',
       routes_to varchar(255) NOT NULL DEFAULT '',
       price varchar(255) NOT NULL DEFAULT '',
       date varchar(255) NOT NULL DEFAULT '',
       time varchar(255) NOT NULL DEFAULT '',
       lastupdated VARCHAR(255) DEFAULT '',
       routeid uuid NOT NULL DEFAULT uuid_generate_v4()
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "routes";

