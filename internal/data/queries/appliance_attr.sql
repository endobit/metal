--
-- CREATE
--

-- name: CreateApplianceAttribute :exec
INSERT INTO attributes (
	entity,
	object,
	name
)
VALUES (
	(SELECT id FROM entities WHERE name = 'appliance'),
	(SELECT id FROM appliances app WHERE app.name = @appliance AND zone = (SELECT id FROM zones z WHERE z.name = @zone)),
	@name
);

--
-- READ
--

-- name: ReadApplianceAttributes :many
SELECT
	a.id,
	z.name AS zone,
	app.name AS appliance,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	appliances app ON a.object = app.id
JOIN
	zones z ON app.zone = z.id
WHERE
	e.name = 'appliance'
ORDER BY
        z.name,
	app.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadApplianceAttributesByAppliance :many
SELECT
	a.id,
	z.name AS zone,
	app.name AS appliance,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	appliances app ON a.object = app.id
JOIN
	zones z ON app.zone = z.id
WHERE
	e.name = 'appliance'
	AND z.name = @zone
	AND app.name = @appliance
ORDER BY
	app.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadApplianceAttributesByZone :many
SELECT
	a.id,
	z.name AS zone,
	app.name AS appliance,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	appliances app ON a.object = app.id
JOIN
	zones z ON app.zone = z.id
WHERE
	e.name = 'appliance'
	AND z.name = @zone
ORDER BY
	z.name,
	app.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadApplianceAttribute :one
SELECT
	a.id,
	z.name AS zone,
	app.name AS appliance,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	appliances app ON a.object = app.id
JOIN
	zones z ON app.zone = z.id
WHERE
	e.name = 'appliance'
	AND z.name = @zone
	AND app.name = @appliance
	AND a.name = @name;

-- name: ReadApplianceAttributesByGlob :many
SELECT
	a.id,
	z.name AS zone,
	app.name AS appliance,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	appliances app ON a.object = app.id
JOIN
	zones z ON app.zone = z.id
WHERE
	e.name = 'appliance'
	AND z.name = @zone
	AND app.name = @appliance
	AND a.name GLOB @glob
ORDER BY
	app.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateApplianceAttributeName :exec
UPDATE
	attributes
SET
	name = @name
WHERE
	entity = (SELECT id FROM entities WHERE name = 'appliance')
	AND object = (
	    SELECT
		id
	    FROM
		appliances app
	    WHERE
	    	app.name = @appliance
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

-- name: UpdateApplianceAttributeValue :exec
UPDATE
	attributes
SET
	value = @value
WHERE
	entity = (SELECT id FROM entities WHERE name = 'appliance')
	AND object = (
	    SELECT
		id
	    FROM
		appliances app
	    WHERE
	    	app.name = @appliance
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

-- name: UpdateApplianceAttributeProtection :exec
UPDATE
	attributes
SET
	is_protected = @is_protected
WHERE
	entity = (SELECT id FROM entities WHERE name = 'appliance')
	AND object = (
	    SELECT
		id
	    FROM
		appliances app
	    WHERE
	    	app.name = @appliance
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

--
-- DELETE
--


