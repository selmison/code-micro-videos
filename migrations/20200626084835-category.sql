
-- +migrate Up
CREATE TABLE categories (
  id uuid NOT NULL PRIMARY KEY,
  name varchar(255) NOT NULL UNIQUE,
  description text,
  is_validated boolean NOT NULL DEFAULT true,
  created_at timestamp,
  updated_at timestamp
);

-- +migrate Down
DROP TABLE categories;
