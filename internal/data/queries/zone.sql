--
-- CREATE
--

-- name: CreateZone :exec
INSERT INTO zones (
	name
)
VALUES (
	@name
);

--
-- READ
--

-- name: ReadZones :many
SELECT
	id,
	name,
	time_zone
FROM
	zones
ORDER BY
	name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadZone :one
SELECT
	id,
	name,
	time_zone
FROM
	zones
WHERE
	name = @name;

-- name: ReadZonesByGlob :many
SELECT
	id,
	name,
	time_zone
FROM
	zones
WHERE
	name GLOB @glob
ORDER BY
	name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateZoneName :exec
UPDATE
	zones
SET
	name = @name
WHERE
	name = @zone;

-- name: UpdateZoneTimeZone :exec
UPDATE
	zones
SET
	time_zone = @time_zone
WHERE
	name = @zone;

--
-- DELETE
--

-- name: DeleteZone :exec
DELETE FROM
	zones
WHERE
	id = @id;

