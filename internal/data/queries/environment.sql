--
-- CREATE
--

-- name: CreateEnvironment :exec
INSERT INTO environments (
	name,
	zone
)
VALUES (
	@name,
	(SELECT id FROM zones z WHERE z.name = @zone)
);


--
-- READ
--

-- name: ReadEnvironments :many
SELECT
	e.id,
	e.name,
	z.name AS zone
FROM
	environments e
JOIN
	zones z ON e.zone = z.id
ORDER BY
	z.name,
	e.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadEnvironment :one
SELECT
	e.id,
	e.name,
	z.name AS zone
FROM
	environments e
JOIN
	zones z ON e.zone = z.id
WHERE
	e.name = @name
	AND z.name = @zone;

-- name: ReadEnvironmentsByGlob :many
SELECT
	e.id,
	e.name,
	z.name AS zone
FROM
	environments e
JOIN
	zones z ON e.zone = z.id
WHERE
	z.name = @zone
	AND e.name GLOB @glob
ORDER BY
	e.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadEnvironmentsByZone :many
SELECT
	e.id,
	e.name,
	z.name AS zone
FROM
	environments e
JOIN
	zones z ON e.zone = z.id
WHERE
	z.name = @zone
ORDER BY
	e.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateEnvironmentName :exec
UPDATE
	environments
SET
	name = @name
WHERE
	environments.name = @environment
	AND zone = (SELECT id FROM zones z WHERE z.name = @zone);

--
-- DELETE
--

-- name: DeleteEnvironment :exec
DELETE FROM
	environments
WHERE
	environments.id = @id;

