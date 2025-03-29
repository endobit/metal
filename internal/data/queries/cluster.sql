--
-- CREATE
--

-- name: CreateCluster :exec
INSERT INTO clusters (
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

-- name: ReadClusters :many
SELECT
	c.id,
	c.name,
	z.name AS zone
FROM
	clusters c
JOIN
	zones z ON c.zone = z.id
ORDER BY
	z.name,
	c.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadCluster :one
SELECT
	c.id,
	c.name,
	z.name AS zone
FROM
	clusters c
JOIN
	zones z ON c.zone = z.id
WHERE
	c.name = @name
	AND z.name = @zone;

-- name: ReadClustersByGlob :many
SELECT
	c.id,
	c.name,
	z.name AS zone
FROM
	clusters c
JOIN
	zones z ON c.zone = z.id
WHERE
	z.name = @zone
	AND c.name GLOB @glob
ORDER BY
	c.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadClustersByZone :many
SELECT
	c.id,
	c.name,
	z.name AS zone
FROM
	clusters c
JOIN
	zones z ON c.zone = z.id
WHERE
	z.name = @zone
ORDER BY
	c.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateClusterName :exec
UPDATE
	clusters
SET
	name = @name
WHERE
	clusters.name = @cluster
	AND zone = (SELECT id FROM zones z WHERE z.name = @zone);

--
-- DELETE
--

-- name: DeleteCluster :exec
DELETE FROM
	clusters
WHERE
	clusters.id = @id;

