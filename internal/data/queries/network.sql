
--
-- CREATE
--

-- name: CreateNetwork :exec
INSERT INTO networks (
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

-- name: ReadNetworks :many
SELECT
	n.id,
	n.name,
	n.address,
	n.gateway,
	n.is_pxe,
	n.mtu,
	z.name AS zone
FROM
	networks n
JOIN
	zones z ON n.zone = z.id
ORDER BY
	z.name,
	n.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadNetwork :one
SELECT
	n.id,
	n.name,
	n.address,
	n.gateway,
	n.is_pxe,
	n.mtu,
	z.name AS zone
FROM
	networks n
JOIN
	zones z ON n.zone = z.id
WHERE
	n.name = @name
	AND n.zone = (SELECT id FROM zones z WHERE z.name = @zone);

-- name: ReadNetworksByZone :many
SELECT
	n.id,
	n.name,
	n.address,
	n.gateway,
	n.is_pxe,
	n.mtu,
	z.name AS zone
FROM
	networks n
JOIN
	zones z ON n.zone = z.id
WHERE
	z.name = @zone
ORDER BY
	n.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

-- name: ReadNetworksByGlob :many
SELECT
	n.id,
	n.name,
	n.address,
	n.gateway,
	n.is_pxe,
	n.mtu,
	z.name AS zone
FROM
	networks n
JOIN
	zones z ON n.zone = z.id
WHERE
	z.name = @zone
	AND n.name GLOB @glob
ORDER BY
	n.name
LIMIT
	COALESCE(NULLIF(@limit, 0), 100) OFFSET COALESCE(@offset, 0);

--
-- UPDATE
--

-- name: UpdateNetworkName :exec
UPDATE
	networks
SET
	name = @name
WHERE
	networks.name = @network
	AND zone = (SELECT id FROM zones z WHERE z.name = @zone);

-- name: UpdateNetworkAddress :exec
UPDATE
	networks
SET
	address = @address
WHERE
	networks.name = @network
	AND zone = (SELECT id FROM zones z WHERE z.name = @zone);

-- name: UpdateNetworkGateway :exec
UPDATE
	networks
SET
	gateway = @gateway
WHERE
	networks.name = @network
	AND zone = (SELECT id FROM zones z WHERE z.name = @zone);

-- name: UpdateNetworkPXE :exec
UPDATE
	networks
SET
	is_pxe = @is_pxe
WHERE
	networks.name = @network
	AND zone = (SELECT id FROM zones z WHERE z.name = @zone);

-- name: UpdateNetworkMTU :exec
UPDATE
	networks
SET
	mtu = @mtu
WHERE
	networks.name = @network
	AND zone = (SELECT id FROM zones z WHERE z.name = @zone);


--
-- DELETE
--

-- name: DeleteNetwork :exec
DELETE FROM
	networks
WHERE
	id = @id;
