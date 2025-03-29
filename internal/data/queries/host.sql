--
-- CREATE
--

-- name: CreateHost :exec
INSERT INTO hosts (
	zone,
	cluster,
	name
)
VALUES (
	(SELECT id FROM zones z WHERE z.name = @zone),
	(SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)),
	@name
);

--
-- READ
--

-- name: ReadHosts :many
SELECT
	h.id,
	h.name,
	mk.name AS make,
	m.name AS model,
	e.name AS environment,
	a.name AS appliance,
	h.location,
	r.name AS rack,
	h.rank,
	h.slot,
	h.type,
	z.name AS zone,
	c.name AS cluster
FROM
	hosts h
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
LEFT JOIN
	models m ON h.model = m.id
LEFT JOIN
	makes mk ON m.make = mk.id
LEFT JOIN
	environments e ON h.environment = e.id
LEFT JOIN
	appliances a ON h.appliance = a.id
LEFT JOIN
	racks r ON h.rack = r.id
ORDER BY
	z.name,
	COALESCE(c.name, ''),
	h.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadHost :one
SELECT
	h.id,
	h.name,
	mk.name AS make,
	m.name AS model,
	e.name AS environment,
	a.name AS appliance,
	h.location,
	r.name AS rack,
	h.rank,
	h.slot,
	h.type,
	z.name AS zone,
	c.name AS cluster
FROM
	hosts h
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
LEFT JOIN
	models m ON h.model = m.id
LEFT JOIN
	makes mk ON m.make = mk.id
LEFT JOIN
	environments e ON h.environment = e.id
LEFT JOIN
	appliances a ON h.appliance = a.id
LEFT JOIN
	racks r ON h.rack = r.id
WHERE
	h.name = @name
	AND (
		(h.zone IS NOT NULL AND z.name = @zone) -- Standalone host
		OR (h.cluster IS NOT NULL AND c.name = @cluster AND z.name = @zone) -- Clustered host
	);

-- name: ReadHostsByGlob :many
SELECT
	h.id,
	h.name,
	mk.name AS make,
	m.name AS model,
	e.name AS environment,
	a.name AS appliance,
	h.location,
	r.name AS rack,
	h.rank,
	h.slot,
	h.type,
	z.name AS zone,
	c.name AS cluster
FROM
	hosts h
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
LEFT JOIN
	models m ON h.model = m.id
LEFT JOIN
	makes mk ON m.make = mk.id
LEFT JOIN
	environments e ON h.environment = e.id
LEFT JOIN
	appliances a ON h.appliance = a.id
LEFT JOIN
	racks r ON h.rack = r.id
WHERE
	h.name GLOB @glob
	AND (
		(h.zone IS NOT NULL AND z.name = @zone) -- Standalone host
		OR (h.cluster IS NOT NULL AND c.name = @cluster AND z.name = @zone) -- Clustered host
	)
ORDER BY
	z.name,
	COALESCE(c.name, ''),
	h.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadHostsByCluster :many
SELECT
	h.id,
	h.name,
	mk.name AS make,
	m.name AS model,
	e.name AS environment,
	a.name AS appliance,
	h.location,
	r.name AS rack,
	h.rank,
	h.slot,
	h.type,
	z.name AS zone,
	c.name AS cluster
FROM
	hosts h
LEFT JOIN
	zones z ON c.zone = z.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	models m ON h.model = m.id
LEFT JOIN
	makes mk ON m.make = mk.id
LEFT JOIN
	environments e ON h.environment = e.id
LEFT JOIN
	appliances a ON h.appliance = a.id
LEFT JOIN
	racks r ON h.rack = r.id
WHERE
	c.name = @cluster
	AND z.name = @zone

ORDER BY
	z.name,
	c.name,
	h.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadHostsByZone :many
SELECT
	h.id,
	h.name,
	mk.name AS make,
	m.name AS model,
	e.name AS environment,
	a.name AS appliance,
	h.location,
	r.name AS rack,
	h.rank,
	h.slot,
	h.type,
	z.name AS zone,
	c.name AS cluster
FROM
	hosts h
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	models m ON h.model = m.id
LEFT JOIN
	makes mk ON m.make = mk.id
LEFT JOIN
	environments e ON h.environment = e.id
LEFT JOIN
	appliances a ON h.appliance = a.id
LEFT JOIN
	racks r ON h.rack = r.id
WHERE
	z.name = @zone
ORDER BY
	COALESCE(c.name, ''),
	h.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);


--
-- UPDATE
--

-- name: UpdateHostName :exec
UPDATE
	hosts
SET
	name = @name
WHERE
	hosts.name = @host
	AND (
		(zone = (SELECT id FROM zones z WHERE z.name = @zone) AND cluster IS NULL)
		OR (cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
	);

-- name: UpdateHostType :exec
UPDATE
	hosts
SET
	type = @type
WHERE
	hosts.name = @host
	AND (
		(zone = (SELECT id FROM zones z WHERE z.name = @zone) AND cluster IS NULL)
		OR (cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
	);


-- name: UpdateHostModel :exec
UPDATE
	hosts
SET
	model = (
		SELECT m.id
		FROM models m
		JOIN makes mk ON m.make = mk.id
		WHERE m.name = @model AND mk.name = @make
	)
WHERE
	hosts.name = @host
	AND (
		(zone = (SELECT id FROM zones z WHERE z.name = @zone) AND cluster IS NULL)
		OR (cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
	);

-- name: UpdateHostEnvironment :exec
UPDATE
	hosts
SET
	environment = (SELECT id FROM environments e WHERE e.name = @environment)
WHERE
	hosts.name = @host
	AND (
		(zone = (SELECT id FROM zones z WHERE z.name = @zone) AND cluster IS NULL)
		OR (cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
	);

-- name: UpdateHostAppliance :exec
UPDATE
	hosts
SET
	appliance = (SELECT id FROM appliances a WHERE a.name = @appliance)
WHERE
	hosts.name = @host
	AND (
		(zone = (SELECT id FROM zones z WHERE z.name = @zone) AND cluster IS NULL)
		OR (cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
	);

-- name: UpdateHostRack :exec
UPDATE
	hosts
SET
	rack = (SELECT id FROM racks r WHERE r.name = @rack)
WHERE
	hosts.name = @host
	AND (
		(zone = (SELECT id FROM zones z WHERE z.name = @zone) AND cluster IS NULL)
		OR (cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
	);

-- name: UpdateHostLocation :exec
UPDATE
	hosts
SET
	location = @location
WHERE
	hosts.name = @host
	AND (
		(zone = (SELECT id FROM zones z WHERE z.name = @zone) AND cluster IS NULL)
		OR (cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
	);

-- name: UpdateHostRank :exec
UPDATE
	hosts
SET
	rank = @rank
WHERE
	hosts.name = @host
	AND (
		(zone = (SELECT id FROM zones z WHERE z.name = @zone) AND cluster IS NULL)
		OR (cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
	);

-- name: UpdateHostSlot :exec
UPDATE
	hosts
SET
	slot = @slot
WHERE
	hosts.name = @host
	AND (
		(zone = (SELECT id FROM zones z WHERE z.name = @zone) AND cluster IS NULL)
		OR (cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND zone = (SELECT id FROM zones z WHERE z.name = @zone)))
	);

--
-- DELETE
--

-- name: DeleteHost :exec
DELETE FROM
	hosts
WHERE
	id = @id;
