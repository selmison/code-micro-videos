-- +migrate Up
CREATE TABLE videos
(
    id            uuid         NOT NULL PRIMARY KEY,
    title         varchar(255) NOT NULL,
    description   text         NOT NULL,
    year_launched smallint     NOT NULL,
    opened        boolean DEFAULT false,
    rating        smallint     NOT NULL,
    duration      smallint     NOT NULL,
    created_at    timestamp,
    updated_at    timestamp,
    deleted_at    timestamp
);

-- +migrate Down
DROP TABLE videos;
