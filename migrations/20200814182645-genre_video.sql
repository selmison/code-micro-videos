-- +migrate Up
CREATE TABLE genre_video
(
    genre_id uuid NOT NULL REFERENCES genres (id) ON DELETE CASCADE,
    video_id uuid NOT NULL REFERENCES videos (id) ON DELETE CASCADE,
    PRIMARY KEY (genre_id, video_id)
);

-- +migrate Down
DROP TABLE genre_video;
