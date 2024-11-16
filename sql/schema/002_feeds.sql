-- +goose Up
CREATE TABLE feeds(
    name TEXT NOT NULL,
    url TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;