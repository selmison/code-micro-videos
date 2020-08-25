-- +migrate Up
CREATE TABLE category_video
(
    category_id uuid NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    video_id    uuid NOT NULL REFERENCES videos (id) ON DELETE CASCADE,
    UNIQUE (category_id, video_id),
    PRIMARY KEY (category_id, video_id)
);

-- +migrate Down
DROP TABLE category_video;
