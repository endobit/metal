--
-- CREATE
--

-- name: CreateGlobalAttribute :exec
INSERT INTO attributes (
	entity,
	name
)
VALUES (
	(SELECT id FROM entities WHERE name = 'global'),
	@name
);

--
-- READ
--

-- name: ReadGlobalAttributes :many
SELECT
	a.id,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'global')
ORDER BY
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);
	

-- name: ReadGlobalAttribute :one
SELECT
	a.id,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'global')
	AND a.name = @name;


-- name: ReadGlobalAttributesByGlob :many
SELECT
	a.id,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'global')
	AND a.name GLOB @glob
ORDER BY
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateGlobalAttributeName :exec
UPDATE
	attributes
SET
	name = @name
WHERE
	entity = (SELECT id FROM entities WHERE name = 'global')
	AND attributes.name = @attr;

-- name: UpdateGlobalAttributeValue :exec
UPDATE
	attributes
SET
	value = @value
WHERE
	entity = (SELECT id FROM entities WHERE name = 'global')
	AND attributes.name = @attr;

-- name: UpdateGlobalAttributeProtection :exec
UPDATE
	attributes
SET
	is_protected = @is_protected
WHERE
	entity = (SELECT id FROM entities WHERE name = 'global')
	AND attributes.name = @attr;

--
-- DELETE
--

