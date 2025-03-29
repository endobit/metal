--
-- CREATE
--

-- name: CreateRack :exec
INSERT INTO racks (
	zone,
	name
)
VALUES (
	(SELECT id FROM zones z WHERE z.name = @zone),
	@name
);


--
-- READ
--

-- name: ReadRacks :many
SELECT
	r.id,
	r.name,
	z.name AS zone
FROM
	racks r
JOIN
	zones z ON r.zone = z.id
ORDER BY
	z.name,
	r.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadRack :one
SELECT
	r.id,
	r.name,
	z.name AS zone
FROM
	racks r
JOIN
	zones z ON r.zone = z.id
WHERE
	z.name = @zone
	AND r.name = @name;

-- name: ReadRacksByGlob :many
SELECT
	r.id,
	r.name,
	z.name AS zone
FROM
	racks r
JOIN
	zones z ON r.zone = z.id
WHERE
	z.name = @zone
	AND r.name GLOB @glob
ORDER BY
	r.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadRacksByZone :many
SELECT
	r.id,
	r.name,
	z.name AS zone
FROM
	racks r
JOIN
	zones z ON r.zone = z.id
WHERE
	z.name = @zone
ORDER BY
	r.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateRackName :exec
UPDATE
	racks
SET
	name = @name
WHERE
	racks.name = @rack
	AND zone = (SELECT id FROM zones z WHERE z.name = @zone);

--
-- DELETE
--

-- name: DeleteRack :exec
DELETE FROM
	racks
WHERE
	racks.id = @id;

