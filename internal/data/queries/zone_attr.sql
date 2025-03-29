--
-- CREATE
--

-- name: CreateZoneAttribute :exec
INSERT INTO attributes (
	entity,
	object,
	name
)
VALUES (
	(SELECT id FROM entities WHERE name = 'zone'),
	(SELECT id FROM zones z WHERE z.name = @zone),
	@name
);

--
-- READ
--

-- name: ReadZoneAttributes :many
SELECT
	a.id,
	z.name AS zone,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	zones z ON a.object = z.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'zone')
ORDER BY
	z.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadZoneAttributesByZone :many
SELECT
	a.id,
	z.name AS zone,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	zones z ON a.object = z.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'zone')
	AND z.name = @zone
ORDER BY
	z.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadZoneAttribute :one
SELECT
	a.id,
	z.name AS zone,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	zones z ON a.object = z.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'zone')
	AND z.name = @zone
	AND a.name = @name;

-- name: ReadZoneAttributesByGlob :many
SELECT
	a.id,
	z.name AS zone,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	zones z ON a.object = z.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'zone')
	AND z.name = @zone
	AND a.name GLOB @glob
ORDER BY
	z.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateZoneAttributeName :exec
UPDATE
	attributes
SET
	name = @name
WHERE
	entity = (SELECT id FROM entities WHERE name = 'zone')
	AND object = (SELECT id FROM zones z WHERE z.name = @zone)
	AND attributes.name = @attr;

-- name: UpdateZoneAttributeValue :exec
UPDATE
	attributes
SET
	value = @value
WHERE
	entity = (SELECT id FROM entities WHERE name = 'zone')
	AND object = (SELECT id FROM zones z WHERE z.name = @zone)
	AND attributes.name = @attr;

-- name: UpdateZoneAttributeProtection :exec
UPDATE
	attributes
SET
	is_protected = @is_protected
WHERE
	entity = (SELECT id FROM entities WHERE name = 'zone')
	AND object = (SELECT id FROM zones z WHERE z.name = @zone)
	AND attributes.name = @attr;

--
-- DELETE
--


