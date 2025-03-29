--
-- CREATE
--

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (
	user,
	role
)
VALUES (
	@user_id,
	@role_id
);


--
-- READ
--

-- name: ReadRolesForUser :many
SELECT
	r.id,
	r.name,
	r.description
FROM
	roles r
JOIN
	user_roles ur ON r.id = ur.role
WHERE
	ur.user = @user_id
ORDER BY
	r.name;



--
-- UPDATE
--

--
-- DELETE
--

-- name: DeleteRoleFromUser :exec
DELETE FROM
	user_roles
WHERE
	user = @user_id
	AND role = @role_id;
