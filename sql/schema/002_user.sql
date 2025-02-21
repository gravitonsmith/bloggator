-- +goose Up
CREATE TABLE users (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT NOT NULL
);

CREATE TABLE feed (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	user_id UUID NOT NULL,
	FOREIGN KEY (user_id)
	REFERENCES users(id)
	ON DELETE CASCADE
);
-- +goose Down
DROP TABLE feed;
DROP TABLE users;
