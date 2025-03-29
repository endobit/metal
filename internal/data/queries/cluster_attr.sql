--
-- CREATE
--

-- name: CreateClusterAttribute :exec
INSERT INTO attributes (
	entity,
	object,
	name
)
VALUES (
	(SELECT id FROM entities WHERE name = 'cluster'),
	(SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)),
	@name
);

--
-- READ
--

-- name: ReadClusterAttributes :many
SELECT
	a.id,
	z.name AS zone,
	c.name AS cluster,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	clusters c ON a.object = c.id
JOIN
	zones z ON c.zone = z.id
WHERE
	e.name = 'cluster'
ORDER BY
        z.name,
	c.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadClusterAttributesByCluster :many
SELECT
	a.id,
	z.name AS zone,
	c.name AS cluster,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	clusters c ON a.object = c.id
JOIN
	zones z ON c.zone = z.id
WHERE
	e.name = 'cluster'
	AND z.name = @zone
	AND c.name = @cluster
ORDER BY
	c.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadClusterAttributesByZone :many
SELECT
	a.id,
	z.name AS zone,
	c.name AS cluster,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	clusters c ON a.object = c.id
JOIN
	zones z ON c.zone = z.id
WHERE
	e.name = 'cluster'
	AND z.name = @zone
ORDER BY
	z.name,
	c.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadClusterAttribute :one
SELECT
	a.id,
	z.name AS zone,
	c.name AS cluster,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	clusters c ON a.object = c.id
JOIN
	zones z ON c.zone = z.id
WHERE
	e.name = 'cluster'
	AND z.name = @zone
	AND c.name = @cluster
	AND a.name = @name;

-- name: ReadClusterAttributesByGlob :many
SELECT
	a.id,
	z.name AS zone,
	c.name AS cluster,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	clusters c ON a.object = c.id
JOIN
	zones z ON c.zone = z.id
WHERE
	e.name = 'cluster'
	AND z.name = @zone
	AND c.name = @cluster
	AND a.name GLOB @glob
ORDER BY
	c.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateClusterAttributeName :exec
UPDATE
	attributes
SET
	name = @name
WHERE
	entity = (SELECT id FROM entities WHERE name = 'cluster')
	AND object = (
	    SELECT
		id
	    FROM
		clusters c
	    WHERE
	    	c.name = @cluster
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

-- name: UpdateClusterAttributeValue :exec
UPDATE
	attributes
SET
	value = @value
WHERE
	entity = (SELECT id FROM entities WHERE name = 'cluster')
	AND object = (
	    SELECT
		id
	    FROM
		clusters c
	    WHERE
	    	c.name = @cluster
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

-- name: UpdateClusterAttributeProtection :exec
UPDATE
	attributes
SET
	is_protected = @is_protected
WHERE
	entity = (SELECT id FROM entities WHERE name = 'cluster')
	AND object = (
	    SELECT
		id
	    FROM
		clusters c
	    WHERE
	    	c.name = @cluster
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

--
-- DELETE
--


