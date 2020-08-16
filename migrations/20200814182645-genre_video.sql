-- +migrate Up
CREATE TABLE genre_video
(
    genre_id uuid NOT NULL REFERENCES genres (id),
    video_id uuid NOT NULL REFERENCES videos (id),
    UNIQUE (genre_id, video_id),
    PRIMARY KEY (genre_id, video_id)
);

-- +migrate Down
DROP TABLE genre_video;
