--
-- TODO: Remove unused queries
--


--
-- CREATE
--

-- name: CreateAttribute :one
INSERT INTO attributes (
	entity,
	object,
	name
)
VALUES (
	(SELECT id FROM entities e WHERE e.name = @entity),
	@object,
	@name
)
RETURNING id;

--
-- READ
--
-- Reading attributes directly from this table is not really useful.


-- name: ReadAttributes :many
SELECT
	a.id,
	a.entity AS entity_id,
	e.name AS entity_name,	-- global, zone, cluster, etc
	a.object,		-- object id, or NULL for global attributes
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
WHERE
	a.entity = (SELECT id FROM entities e WHERE e.name = @entity)
ORDER BY
	a.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateAttribute :exec
UPDATE
	attributes
SET
	value = @value,
	is_protected = @is_protected
WHERE
	entity = (SELECT id FROM entities e WHERE e.name = @entity)
	AND object = @object
	AND attributes.name = @name;

-- name: UpdateAttributeName :exec
UPDATE
	attributes
SET
	name = @new_name
WHERE
	entity = (SELECT id FROM entities e WHERE e.name = @entity)
	AND object = @object
	AND attributes.name = @name;

-- name: UpdateAttributeValue :exec
UPDATE
	attributes
SET
	value = @value
WHERE
	entity = (SELECT id FROM entities e WHERE e.name = @entity)
	AND object = @object
	AND attributes.name = @name;

-- name: UpdateAttributeProtection :exec
UPDATE
	attributes
SET
	is_protected = @is_protected
WHERE
	entity = (SELECT id FROM entities e WHERE e.name = @entity)
	AND object = @object
	AND attributes.name = @name;

--
-- DELETE
--
-- Note there is no cascade delete on object deletion.

-- name: DeleteAttribute :exec
DELETE FROM
	attributes
WHERE
	id = @id;

