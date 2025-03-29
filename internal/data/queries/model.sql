--
-- CREATE
--

-- name: CreateModel :exec
INSERT INTO models (
	make,
	name
)
VALUES (
	(SELECT id FROM makes m WHERE m.name = @make),
	@name
);

--
-- READ
--

-- name: ReadModels :many
SELECT
	m.id,
	mk.name AS make,
	m.name,
	m.architecture
FROM
	models m
JOIN
	makes mk ON m.make = mk.id
ORDER BY
	mk.name,
	m.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadModelsByMake :many
SELECT
	m.id,
	mk.name AS make,
	m.name,
	m.architecture
FROM
	models m
JOIN
	makes mk ON m.make = mk.id
WHERE
	mk.name = @make
ORDER BY
	mk.name,
	m.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);


-- name: ReadModel :one
SELECT
	m.id,
	mk.name AS make,
	m.name,
	m.architecture
FROM
	models m
JOIN
	makes mk ON m.make = mk.id
WHERE
	mk.name = @make
	AND m.name = @name;

-- name: ReadModelsByGlob :many
SELECT
	m.id,
	mk.name AS make,
	m.name,
	m.architecture
FROM
	models m
JOIN
	makes mk ON m.make = mk.id
WHERE
	m.name GLOB @glob
	AND mk.name = @make
ORDER BY
	mk.name,
	m.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);


--
-- UPDATE
--

-- name: UpdateModelName :exec
UPDATE
	models
SET
	name = @name
WHERE
	models.name = @model
	AND make = (SELECT id FROM makes m WHERE m.name = @make);

-- name: UpdateModelArchitecture :exec
UPDATE
	models
SET
	architecture = @architecture
WHERE
	models.name = @model
	AND make = (SELECT id FROM makes m WHERE m.name = @make);


--
-- DELETE
--

-- name: DeleteModel :exec
DELETE FROM
	models
WHERE
	id = @id;
