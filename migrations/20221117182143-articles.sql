
-- +migrate Up
CREATE TABLE IF NOT EXISTS articles (
    id SERIAL,
    category_id INT,
    title varchar(250),
    description Text,
    content Text,
    PRIMARY KEY (id),
    CONSTRAINT fk_category
        FOREIGN KEY(category_id)
            REFERENCES categories(id),
    CONSTRAINT title_unique UNIQUE (title)
);

-- +migrate Down
DROP TABLE IF EXISTS articles;
