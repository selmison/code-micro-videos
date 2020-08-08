-- +migrate Up
CREATE TABLE genres
(
    id           uuid         NOT NULL PRIMARY KEY,
    name         varchar(255) NOT NULL UNIQUE,
    is_validated boolean      NOT NULL DEFAULT true,
    created_at   timestamp,
    updated_at   timestamp,
    deleted_at   timestamp
);

-- +migrate Down
DROP TABLE genres;
