-- +goose Up
CREATE TABLE posts (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	title VARCHAR(255) NOT NULL,
	url VARCHAR(255) NOT NULL UNIQUE,
	description TEXT,
	feed_id UUID NOT NULL,
	published_at TIMESTAMP,
	FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);



-- +goose Down
DROP TABLE posts;
