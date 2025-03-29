--
-- CREATE
--

-- name: CreateHostInterface :exec
INSERT INTO host_interfaces (
	host,
	name
)
VALUES (
	(
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	),
	@name
);

--
-- READ
--

-- name: ReadHostInterfaces :many
SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
ORDER BY
	z.name,
	COALESCE(c.name, ''),
	h.name,
	hi.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadHostInterfacesByHost :many
SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
WHERE
	h.name = @host
	AND (
		(h.zone IS NOT NULL AND z.name = @zone)
		OR (h.cluster IS NOT NULL AND c.name = @cluster AND z.name = @zone)
	)
ORDER BY
	z.name,
	COALESCE(c.name, ''),
	h.name,
	hi.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadHostInterfacesByCluster :many
SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
WHERE
	c.name = @cluster
	AND
	z.name = @zone
ORDER BY
	z.name,
	c.name,
	h.name,
	hi.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadHostInterfacesByZone :many
SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
WHERE
	z.name = @zone
ORDER BY
	COALESCE(c.name, ''),
	h.name,
	hi.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadHostInterface :one
SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
WHERE
	h.name = @host
	AND (
		(h.zone IS NOT NULL AND z.name = @zone)
		OR (h.cluster IS NOT NULL AND c.name = @cluster AND z.name = @zone)
	)
	AND hi.name = @name;


-- name: ReadHostInterfacesByGlob :many
SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
LEFT JOIN
	clusters c ON h.cluster = c.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
WHERE
	h.name = @host
	AND (
		(h.zone IS NOT NULL AND z.name = @zone)
		OR (h.cluster IS NOT NULL AND c.name = @cluster AND z.name = @zone)
	)
	AND hi.name GLOB @glob
ORDER BY
        z.name,
	COALESCE(c.name, ''),
	h.name,
	hi.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateHostInterfaceName :exec
UPDATE
	host_interfaces
SET
	name = @name
WHERE
	host_interfaces.name = @interface
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);


-- name: UpdateHostInterfaceIP :exec
UPDATE
	host_interfaces
SET
	ip = @ip
WHERE
	host_interfaces.name = @interface
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

-- name: UpdateHostInterfaceMAC :exec
UPDATE
	host_interfaces
SET
	mac = @mac
WHERE
	host_interfaces.name = @interface
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

-- name: UpdateHostInterfaceNetmask :exec
UPDATE
	host_interfaces
SET
	netmask = @mask
WHERE
	host_interfaces.name = @interface
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

-- name: UpdateHostInterfaceDHCP :exec
UPDATE
	host_interfaces
SET
	is_dhcp = @dhcp
WHERE
	host_interfaces.name = @interface
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

-- name: UpdateHostInterfacePXE :exec
UPDATE
	host_interfaces
SET
	is_pxe = @pxe
WHERE
	host_interfaces.name = @interface
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

-- name: UpdateHostInterfaceManagement :exec
UPDATE
	host_interfaces
SET
	is_management = @management
WHERE
	host_interfaces.name = @interface
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

-- name: UpdateHostInterfaceType :exec
UPDATE
	host_interfaces
SET
	type = @type
WHERE
	host_interfaces.name = @interface
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

-- name: UpdateHostInterfaceBondMode :exec
UPDATE
	host_interfaces
SET
	bond_mode = @bond
WHERE
	host_interfaces.name = @interface
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

-- name: UpdateHostInterfaceNetwork :exec
UPDATE
	host_interfaces
SET
	network = (
		SELECT id
		FROM networks n
		WHERE n.name = @network
		AND n.zone = (SELECT id FROM zones z WHERE z.name = @zone)
	)
WHERE
	host_interfaces.name = @interface
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);

-- name: UpdateHostInterfaceMaster :exec
UPDATE
	host_interfaces
SET
	master = (
		SELECT id
		FROM host_interfaces hi
		WHERE hi.name = @master
		AND host = (
			SELECT id
			FROM hosts h
			WHERE h.name = @host
			AND (
				(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
			)
		)
	)
WHERE
	host_interfaces.name = @interface
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = @host
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = @zone) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = @cluster AND h.zone = (SELECT id FROM zones z WHERE z.name = @zone)))
		)
	);




--
-- DELETE
--

-- name: DeleteHostInterface :exec
DELETE FROM
	host_interfaces
WHERE
	id = @id;
