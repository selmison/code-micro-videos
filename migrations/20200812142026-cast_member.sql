-- +migrate Up
CREATE TABLE cast_members
(
    id         uuid         NOT NULL PRIMARY KEY,
    name       varchar(255) NOT NULL,
    type       smallint     NOT NULL,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

-- +migrate Down
DROP TABLE cast_members;
