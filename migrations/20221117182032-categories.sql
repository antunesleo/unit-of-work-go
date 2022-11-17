
-- +migrate Up
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL,
    name varchar(250),
    PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS categories;
