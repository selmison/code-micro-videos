-- +migrate Up
CREATE TABLE category_video
(
    category_id uuid NOT NULL REFERENCES categories (id),
    video_id    uuid NOT NULL REFERENCES videos (id),
    UNIQUE (category_id, video_id),
    PRIMARY KEY (category_id, video_id)
);

-- +migrate Down
DROP TABLE category_video;
