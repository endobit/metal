--
-- CREATE
--

-- name: CreateRole :exec
INSERT INTO roles (
	name
)
VALUES (
	@name
);



--
-- READ
--

-- name: ReadRoles :many
SELECT
	id,
	name,
	description
FROM
	roles
ORDER BY
	name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);


-- name: ReadRole :one
SELECT
	id,
	name,
	description
FROM
	roles
WHERE
	name = @role;

-- name: ReadRolesByGlob :many
SELECT
	id,
	name,
	description
FROM
	roles
WHERE
	name GLOB @glob
ORDER BY
	name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);



--
-- UPDATE
--

-- name: UpdateRoleName :exec
UPDATE
	roles
SET
	name = @name
WHERE
	name = @role;

-- name: UpdateRoleDescription :exec
UPDATE
	roles
SET
	description = @description
WHERE
	name = @role;

--
-- DELETE
--

-- name: DeleteRole :exec
DELETE FROM
	roles
WHERE
	id = @id;


