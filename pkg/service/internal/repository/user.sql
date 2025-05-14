-- name: create
INSERT INTO users (id) VALUES ($1);

-- name: exists
SELECT EXISTS (SELECT 1 FROM users WHERE id = $1);
