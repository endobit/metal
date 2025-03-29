--
-- CREATE
--

-- name: CreateUser :exec
INSERT INTO users (
	name
)
VALUES (
	@name
);


--
-- READ
--

-- name: ReadUsers :many
SELECT
	id,
	name,
	email
FROM
	users
ORDER BY
	name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);


-- name: ReadUser :one
SELECT
	id,
	name,
	password_hash,
	email
FROM
	users
WHERE
	name = @user;

-- name: ReadUsersByGlob :many
SELECT
	id,
	name,
	password_hash,
	email
FROM
	users
WHERE
	name GLOB @glob
ORDER BY
	name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);


--
-- UPDATE
--

-- name: UpdateUserEmail :exec
UPDATE
	users
SET
	email = @email
WHERE
	name = @user;

-- name: UpdateUserPassword :exec
UPDATE
	users
SET
	password_hash = @password_hash
WHERE
	name = @user;


--
-- DELETE
--

-- name: DeleteUser :exec
DELETE FROM
	users
WHERE
	id = @id;

