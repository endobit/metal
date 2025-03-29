--
-- CREATE
--

-- name: CreateMake :exec
INSERT INTO makes (
	name
)
VALUES (
	@name
);

--
-- READ
--

-- name: ReadMakes :many
SELECT
	id,
	name
FROM
	makes
ORDER BY
	name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadMake :one
SELECT
	id,
	name
FROM
	makes
WHERE
	name = @name;

-- name: ReadMakesByGlob :many
SELECT
	id,
	name
FROM
	makes
WHERE
	name GLOB @glob
ORDER BY
	name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateMakeName :exec
UPDATE
	makes
SET
	name = @name
WHERE
	name = @make;

--
-- DELETE
--

-- name: DeleteMake :exec
DELETE FROM
	makes
WHERE
	id = @id;
