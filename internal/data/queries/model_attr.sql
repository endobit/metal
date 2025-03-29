--
-- CREATE
--

-- name: CreateModelAttribute :exec
INSERT INTO attributes (
	entity,
	object,
	name
)
VALUES (
	(SELECT id FROM entities WHERE name = 'model'),
	(SELECT id FROM models m, makes mk WHERE m.name = @model AND m.make = mk.id AND mk.name = @make),
	@name
);

--
-- READ
--

-- name: ReadModelAttributes :many
SELECT
	a.id,
	m.name AS model,
	mk.name AS make,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	models m ON a.object = m.id
JOIN
	makes mk ON m.make = mk.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'model')
ORDER BY
        mk.name,
	m.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadModelAttributesByMake :many
SELECT
	a.id,
	m.name AS model,
	mk.name AS make,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	models m ON a.object = m.id
JOIN
	makes mk ON m.make = mk.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'model')
	AND mk.name = @make
ORDER BY
        mk.name,
	m.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadModelAttributesByMakeModel :many
SELECT
	a.id,
	m.name AS model,
	mk.name AS make,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	models m ON a.object = m.id
JOIN
	makes mk ON m.make = mk.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'model')
	AND mk.name = @make
	AND m.name = @model
ORDER BY
        mk.name,
	m.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadModelAttribute :one
SELECT
	a.id,
	m.name AS model,
	mk.name AS make,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	models m ON a.object = m.id
JOIN
	makes mk ON m.make = mk.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'model')
	AND mk.name = @make
	AND m.name = @model
	AND a.name = @name;

-- name: ReadModelAttributesByGlob :many
SELECT
	a.id,
	m.name AS model,
	mk.name AS make,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	models m ON a.object = m.id
JOIN
	makes mk ON m.make = mk.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'model')
	AND mk.name = @make
	AND m.name = @model
	AND a.name GLOB @glob
ORDER BY
        mk.name,
	m.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateModelAttributeName :exec
UPDATE
	attributes
SET
	name = @name
WHERE
	entity = (SELECT id FROM entities WHERE name = 'model')
	AND object = (SELECT id FROM models m, makes mk WHERE m.name = @model AND m.make = mk.id AND mk.name = @make)
	AND attributes.name = @attr;

-- name: UpdateModelAttributeValue :exec
UPDATE
	attributes
SET
	value = @value
WHERE
	entity = (SELECT id FROM entities WHERE name = 'model')
	AND object = (SELECT id FROM models m, makes mk WHERE m.name = @model AND m.make = mk.id AND mk.name = @make)
	AND attributes.name = @attr;

-- name: UpdateModelAttributeProtection :exec
UPDATE
	attributes
SET
	is_protected = @is_protected
WHERE
	entity = (SELECT id FROM entities WHERE name = 'model')
	AND object = (SELECT id FROM models m, makes mk WHERE m.name = @model AND m.make = mk.id AND mk.name = @make)
	AND attributes.name = @attr;

--
-- DELETE
--


