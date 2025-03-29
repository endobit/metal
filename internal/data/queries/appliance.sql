--
-- CREATE
--

-- name: CreateAppliance :exec
INSERT INTO appliances (
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

-- name: ReadAppliances :many
SELECT
	a.id,
	a.name,
	z.name AS zone
FROM
	appliances a
JOIN
	zones z ON a.zone = z.id
ORDER BY
	z.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadAppliance :one
SELECT
	a.id,
	a.name,
	z.name AS zone
FROM
	appliances a
JOIN
	zones z ON a.zone = z.id
WHERE
	a.name = @name
	AND z.name = @zone;

-- name: ReadAppliancesByGlob :many
SELECT
	a.id,
	a.name,
	z.name AS zone
FROM
	appliances a
JOIN
	zones z ON a.zone = z.id
WHERE
	z.name = @zone
	AND a.name GLOB @glob
ORDER BY
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadAppliancesByZone :many
SELECT
	a.id,
	a.name,
	z.name AS zone
FROM
	appliances a
JOIN
	zones z ON a.zone = z.id
WHERE
	z.name = @zone
ORDER BY
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateApplianceName :exec
UPDATE
	appliances
SET
	name = @name
WHERE
	appliances.name = @appliance
	AND zone = (SELECT id FROM zones z WHERE z.name = @zone);

--
-- DELETE
--

-- name: DeleteAppliance :exec
DELETE FROM
	appliances
WHERE
	appliances.id = @id;

