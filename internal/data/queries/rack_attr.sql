--
-- CREATE
--

-- name: CreateRackAttribute :exec
INSERT INTO attributes (
	entity,
	object,
	name
)
VALUES (
	(SELECT id FROM entities WHERE name = 'rack'),
	(SELECT id FROM racks r WHERE r.name = @rack AND zone = (SELECT id FROM zones z WHERE z.name = @zone)),
	@name
);

--
-- READ
--

-- name: ReadRackAttributes :many
SELECT
	a.id,
	z.name AS zone,
	r.name AS rack,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	racks r ON a.object = r.id
JOIN
	zones z ON r.zone = z.id
WHERE
	e.name = 'rack'
ORDER BY
        z.name,
	r.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadRackAttributesByRack :many
SELECT
	a.id,
	z.name AS zone,
	r.name AS rack,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	racks r ON a.object = r.id
JOIN
	zones z ON r.zone = z.id
WHERE
	e.name = 'rack'
	AND z.name = @zone
	AND r.name = @rack
ORDER BY
	r.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadRackAttributesByZone :many
SELECT
	a.id,
	z.name AS zone,
	r.name AS rack,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	racks r ON a.object = r.id
JOIN
	zones z ON r.zone = z.id
WHERE
	e.name = 'rack'
	AND z.name = @zone
ORDER BY
	z.name,
	r.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadRackAttribute :one
SELECT
	a.id,
	z.name AS zone,
	r.name AS rack,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	racks r ON a.object = r.id
JOIN
	zones z ON r.zone = z.id
WHERE
	e.name = 'rack'
	AND z.name = @zone
	AND r.name = @rack
	AND a.name = @name;

-- name: ReadRackAttributesByGlob :many
SELECT
	a.id,
	z.name AS zone,
	r.name AS rack,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	racks r ON a.object = r.id
JOIN
	zones z ON r.zone = z.id
WHERE
	e.name = 'rack'
	AND z.name = @zone
	AND r.name = @rack
	AND a.name GLOB @glob
ORDER BY
	r.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateRackAttributeName :exec
UPDATE
	attributes
SET
	name = @name
WHERE
	entity = (SELECT id FROM entities WHERE name = 'rack')
	AND object = (
	    SELECT
		id
	    FROM
		racks r
	    WHERE
	    	r.name = @rack
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

-- name: UpdateRackAttributeValue :exec
UPDATE
	attributes
SET
	value = @value
WHERE
	entity = (SELECT id FROM entities WHERE name = 'rack')
	AND object = (
	    SELECT
		id
	    FROM
		racks r
	    WHERE
	    	r.name = @rack
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

-- name: UpdateRackAttributeProtection :exec
UPDATE
	attributes
SET
	is_protected = @is_protected
WHERE
	entity = (SELECT id FROM entities WHERE name = 'rack')
	AND object = (
	    SELECT
		id
	    FROM
		racks r
	    WHERE
	    	r.name = @rack
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

--
-- DELETE
--


