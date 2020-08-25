-- +migrate Up
CREATE TABLE category_genre
(
    category_id uuid NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    genre_id    uuid NOT NULL REFERENCES genres (id) ON DELETE CASCADE,
    PRIMARY KEY (category_id, genre_id)
);

-- +migrate Down
DROP TABLE category_genre;
