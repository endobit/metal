--
-- CREATE
--

-- name: CreateEnvironmentAttribute :exec
INSERT INTO attributes (
	entity,
	object,
	name
)
VALUES (
	(SELECT id FROM entities WHERE name = 'environment'),
	(SELECT id FROM environments env WHERE env.name = @environment AND zone = (SELECT id FROM zones z WHERE z.name = @zone)),
	@name
);

--
-- READ
--

-- name: ReadEnvironmentAttributes :many
SELECT
	a.id,
	z.name AS zone,
	env.name AS environment,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	environments env ON a.object = env.id
JOIN
	zones z ON env.zone = z.id
WHERE
	e.name = 'environment'
ORDER BY
        z.name,
	env.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadEnvironmentAttributesByEnvironment :many
SELECT
	a.id,
	z.name AS zone,
	env.name AS environment,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	environments env ON a.object = env.id
JOIN
	zones z ON env.zone = z.id
WHERE
	e.name = 'environment'
	AND z.name = @zone
	AND env.name = @environment
ORDER BY
	env.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadEnvironmentAttributesByZone :many
SELECT
	a.id,
	z.name AS zone,
	env.name AS environment,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	environments env ON a.object = env.id
JOIN
	zones z ON env.zone = z.id
WHERE
	e.name = 'environment'
	AND z.name = @zone
ORDER BY
	z.name,
	env.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadEnvironmentAttribute :one
SELECT
	a.id,
	z.name AS zone,
	env.name AS environment,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	environments env ON a.object = env.id
JOIN
	zones z ON env.zone = z.id
WHERE
	e.name = 'environment'
	AND z.name = @zone
	AND env.name = @environment
	AND a.name = @attr;

-- name: ReadEnvironmentAttributesByGlob :many
SELECT
	a.id,
	z.name AS zone,
	env.name AS environment,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	environments env ON a.object = env.id
JOIN
	zones z ON env.zone = z.id
WHERE
	e.name = 'environment'
	AND z.name = @zone
	AND env.name = @environment
	AND a.name GLOB @glob
ORDER BY
	env.name,
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateEnvironmentAttributeName :exec
UPDATE
	attributes
SET
	name = @name
WHERE
	entity = (SELECT id FROM entities WHERE name = 'environment')
	AND object = (
	    SELECT
		id
	    FROM
		environments env
	    WHERE
	    	env.name = @environment
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

-- name: UpdateEnvironmentAttributeValue :exec
UPDATE
	attributes
SET
	value = @value
WHERE
	entity = (SELECT id FROM entities WHERE name = 'environment')
	AND object = (
	    SELECT
		id
	    FROM
		environments env
	    WHERE
	    	env.name = @environment
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

-- name: UpdateEnvironmentAttributeProtection :exec
UPDATE
	attributes
SET
	is_protected = @is_protected
WHERE
	entity = (SELECT id FROM entities WHERE name = 'environment')
	AND object = (
	    SELECT
		id
	    FROM
		environments env
	    WHERE
	    	env.name = @environment
	    	AND zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
	AND attributes.name = @attr;

--
-- DELETE
--


