--
-- CREATE
--

-- name: CreateHostAttribute :exec
INSERT INTO attributes (
	entity,
	object,
	name
)
VALUES (
	(SELECT id FROM entities WHERE name = 'host'),
	(SELECT id FROM hosts h WHERE h.name = @host AND (
		(zone = (SELECT id FROM zones z WHERE z.name = @zone) AND cluster IS NULL)
		OR (cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
	)),
	@name
);

--
-- READ
--

-- name: ReadHostAttributes :many
SELECT
	a.id,
	a.name,
	a.value,
	a.is_protected,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	hosts h ON a.object = h.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
WHERE
	e.name = 'host'
ORDER BY
	z.name,
	COALESCE(c.name, ''),
	h.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);


-- name: ReadHostAttributesByHost :many
SELECT
	a.id,
	a.name,
	a.value,
	a.is_protected,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	hosts h ON a.object = h.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
WHERE
	e.name = 'host'
	AND h.name = @host
	AND (
		(h.zone IS NOT NULL AND z.name = @zone)
		OR (h.cluster IS NOT NULL AND c.name = @cluster AND z.name = @zone)
	)
ORDER BY
	z.name,
	COALESCE(c.name, ''),
	h.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);


-- name: ReadHostAttributesByCluster :many
SELECT
	a.id,
	a.name,
	a.value,
	a.is_protected,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	hosts h ON a.object = h.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
WHERE
	e.name = 'host'
	AND c.name = @cluster
	AND z.name = @zone

ORDER BY
	z.name,
	c.name,
	h.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);


-- name: ReadHostAttribute :one
SELECT
	a.id,
	a.name,
	a.value,
	a.is_protected,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	hosts h ON a.object = h.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
WHERE
	e.name = 'host'
	AND h.name = @host
	AND (
		(h.zone IS NOT NULL AND z.name = @zone)
		OR (h.cluster IS NOT NULL AND c.name = @cluster AND z.name = @zone)
	)
	AND a.name = @name;

-- name: ReadHostAttributesByGlob :many
SELECT
	a.id,
	a.name,
	a.value,
	a.is_protected,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	hosts h ON a.object = h.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
WHERE
	e.name = 'host'
	AND h.name = @host
	AND (
		(h.zone IS NOT NULL AND z.name = @zone)
		OR (h.cluster IS NOT NULL AND c.name = @cluster AND z.name = @zone)
	)
	AND a.name GLOB @glob
ORDER BY
	z.name,
	COALESCE(c.name, ''),
	h.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);


-- name: ReadHostAttributesByZone :many
SELECT
	a.id,
	a.name,
	a.value,
	a.is_protected,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	hosts h ON a.object = h.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
WHERE
	e.name = 'host'
	AND z.name = @zone
ORDER BY
	z.name,
	COALESCE(c.name, ''),
	h.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateHostAttributeName :exec
UPDATE
	attributes
SET
	name = @name
WHERE
	attributes.name = @attr
	AND entity = (SELECT id FROM entities WHERE name = 'host')
	AND object = (
		SELECT
			id
		FROM
			hosts h
		WHERE
			h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

-- name: UpdateHostAttributeValue :exec
UPDATE
	attributes
SET
	value = @value
WHERE
	attributes.name = @attr
	AND entity = (SELECT id FROM entities WHERE name = 'host')
	AND object = (
		SELECT
			id
		FROM
			hosts h
		WHERE
			h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

-- name: UpdateHostAttributeProtection :exec
UPDATE
	attributes
SET
	is_protected = @is_protected
WHERE
	attributes.name = @attr
	AND entity = (SELECT id FROM entities WHERE name = 'host')
	AND object = (
		SELECT
			id
		FROM
			hosts h
		WHERE
			h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

--
-- DELETE
--
