-- +goose Up

CREATE TABLE entities (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the entity
	name		TEXT NOT NULL UNIQUE	-- Type of the entity (e.g., 'zone', 'cluster', 'model', 'global')
);

-- The hierarchy of entities is determined by sorting on the id column.

INSERT INTO entities (name) VALUES ('global');
INSERT INTO entities (name) VALUES ('model');
INSERT INTO entities (name) VALUES ('zone');
INSERT INTO entities (name) VALUES ('appliance');
INSERT INTO entities (name) VALUES ('rack');
INSERT INTO entities (name) VALUES ('cluster');
INSERT INTO entities (name) VALUES ('environment');
INSERT INTO entities (name) VALUES ('host');
INSERT INTO entities (name) VALUES ('switch');
INSERT INTO entities (name) VALUES ('bmc');

CREATE TABLE makes (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the make
	name		TEXT NOT NULL UNIQUE	-- Name of the make (must be unique)
);

CREATE TABLE models (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the model
	make		INTEGER NOT NULL,	-- Foreign key to the makes table
	name		TEXT NOT NULL,		-- Model name
	architecture	TEXT,			-- Architecture type
	UNIQUE (make, name),			-- Ensure make and model combination is unique
	FOREIGN KEY (make) REFERENCES makes(id) ON DELETE CASCADE	-- Cascade deletion when a make is deleted
);

CREATE TABLE oses (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the OS
	name		TEXT NOT NULL UNIQUE	-- Name of the OS
);

CREATE TABLE roles (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the role
	name		TEXT NOT NULL UNIQUE,	-- Name of the role (must be unique)
	description	TEXT			-- Description of the role
);

CREATE TABLE users (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the user
	name		TEXT NOT NULL UNIQUE,	-- Username of the user (must be unique)
	password_hash	TEXT,			-- Hashed password for secure storage
	email		TEXT UNIQUE		-- Email address of the user (optional, but must be unique if provided)
);

CREATE TABLE zones (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the zone
	name		TEXT NOT NULL UNIQUE,	-- Name of the zone (must be unique)
	time_zone	TEXT			-- Time zone associated with the zone
);




CREATE TABLE appliances (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the appliance
	name		TEXT NOT NULL,		-- Name of the appliance
	zone		INTEGER NOT NULL,	-- Foreign key to the zones table
	UNIQUE (name, zone),			-- Combination of name and zone must be unique
	FOREIGN KEY (zone) REFERENCES zones(id) ON DELETE CASCADE
);

CREATE TABLE attributes (
	id		INTEGER PRIMARY KEY,		-- Unique identifier for the attribute
	entity		INTEGER NOT NULL,		-- Foreign key to the entities table
	object		INTEGER,			-- ID of the associated object (NULL for global attributes)
	name		TEXT NOT NULL,			-- Attribute key
	value		TEXT,				-- Attribute value
	is_protected	INTEGER CHECK (is_protected IN (0, 1)),	-- Whether the attribute is protected
	UNIQUE (entity, object, name),		-- Ensure uniqueness for (entity, object, key) combinations
	FOREIGN KEY (entity) REFERENCES entities(id) ON DELETE CASCADE
);

CREATE TABLE clusters (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the cluster
	name		TEXT NOT NULL,		-- Name of the cluster
	zone		INTEGER NOT NULL,	-- Foreign key to the zones table
	UNIQUE (name, zone),			-- Combination of name and zone must be unique
	FOREIGN KEY (zone) REFERENCES zones(id) ON DELETE CASCADE
);

CREATE TABLE environments (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the environment
	name		TEXT NOT NULL UNIQUE,	-- Name of the environment
	zone		INTEGER NOT NULL,	-- Foreign key to the zones table
	UNIQUE (name, zone),			-- Combination of name and zone must be unique
	FOREIGN KEY (zone) REFERENCES zones(id) ON DELETE CASCADE
);

CREATE TABLE networks (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the network
	name		TEXT NOT NULL,		-- Name of the network
	address		TEXT,			-- Network address
	gateway		TEXT,			-- Network gateway
	is_pxe		INTEGER NOT NULL DEFAULT 0 CHECK (is_pxe IN (0, 1)),	-- Is PXE enabled
	mtu		INTEGER NOT NULL DEFAULT 1500,	-- Maximum transmission unit
	zone		INTEGER NOT NULL,	-- Foreign key to the zones table
	UNIQUE (name, zone),			-- Combination of name and zone must be unique
	FOREIGN KEY (zone) REFERENCES zones(id) ON DELETE CASCADE
);

CREATE TABLE stacks (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the stack
	name		TEXT NOT NULL UNIQUE,	-- Name of the stack
	os		INTEGER NOT NULL,	-- Foreign key to the oses table
	os_flavor	TEXT,			-- OS flavor
	os_version	TEXT,			-- OS version
	FOREIGN KEY (os) REFERENCES oses(id) ON DELETE CASCADE
);

CREATE TABLE racks (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the rack
	name		TEXT NOT NULL,		-- Name of the rack
	zone		INTEGER NOT NULL,	-- Foreign key to the zones table
	UNIQUE (name, zone),			-- Combination of name and zone must be unique
	FOREIGN KEY (zone) REFERENCES zones(id) ON DELETE CASCADE
);

CREATE TABLE user_roles (
	id		INTEGER PRIMARY KEY,		-- Unique identifier for the mapping
	user		INTEGER NOT NULL,		-- Foreign key to the users table
	role		INTEGER NOT NULL,		-- Foreign key to the roles table
	created_at	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,	-- Timestamp of when the mapping was created
	FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (role) REFERENCES roles(id) ON DELETE CASCADE,
	UNIQUE (user, role)				-- Ensure that each user can have a role only once
);


CREATE TABLE device (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the device
	name		TEXT NOT NULL,		-- Name of the device
	model		INTEGER,		-- Foreign key to the models table
	location	TEXT,			-- Location of the device (e.g. office)
	network		INTEGER,		-- Foreign key to the networks table
	zone		INTEGER,		-- Foreign key to the zones table
	FOREIGN KEY (model) REFERENCES models(id),
	FOREIGN KEY (zone) REFERENCES zones(id) ON DELETE CASCADE
);

CREATE TABLE device_interfaces (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the host interface
	name		TEXT NOT NULL,		-- Name of the interface
	ip		TEXT,			-- IP address of the interface
	mac		TEXT,			-- MAC address of the interface
	is_dhcp		INTEGER NOT NULL DEFAULT 0 CHECK (is_dhcp IN (0, 1)),	-- Whether DHCP is enabled
	device		INTEGER NOT NULL,	-- Foreign key to the hosts table
	network		INTEGER,		-- Foreign key to the networks table
	UNIQUE (name, device),			-- Combination of name and device must be unique
	FOREIGN KEY (network) REFERENCES networks(id),
	FOREIGN KEY (device) REFERENCES device(id) ON DELETE CASCADE
);

CREATE TABLE hosts (
	id		INTEGER PRIMARY KEY,		-- Unique identifier for the host
	name		TEXT NOT NULL,			-- Name of the host
	model		INTEGER,			-- Foreign key to the models table
	environment	INTEGER,			-- Foreign key to the environments table
	appliance	INTEGER,			-- Foreign key to the appliances table
	location	TEXT,				-- Location of the host (e.g. office)
	type            TEXT,                           -- Type of the host (e.g., physical, virtual)
	rack		INTEGER,			-- Foreign key to the racks table
	rank		INTEGER,			-- Rank of the host within the rack
	slot		INTEGER,			-- Slot number of the host within the rack
	cluster		INTEGER,			-- Foreign key to the clusters table (if clustered)
	zone		INTEGER,			-- Foreign key to the zones table (if standalone)
	CHECK (cluster IS NOT NULL OR zone IS NOT NULL),
	FOREIGN KEY (model) REFERENCES models(id),
	FOREIGN KEY (environment) REFERENCES environments(id),
	FOREIGN KEY (appliance) REFERENCES appliances(id),
	FOREIGN KEY (rack) REFERENCES racks(id),
	FOREIGN KEY (cluster) REFERENCES clusters(id) ON DELETE CASCADE,
	FOREIGN KEY (zone) REFERENCES zones(id) ON DELETE CASCADE
);

-- Partial index for standalone hosts (where cluster is NULL)
CREATE UNIQUE INDEX unique_standalone_hosts
ON hosts (name, zone)
WHERE cluster IS NULL;

-- Partial index for clustered hosts (where zone is NULL)
CREATE UNIQUE INDEX unique_clustered_hosts
ON hosts (name, cluster)
WHERE zone IS NULL;

CREATE TABLE host_interfaces (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the host interface
	name		TEXT NOT NULL,		-- Name of the interface
	ip		TEXT,			-- IP address of the interface
	mac		TEXT,			-- MAC address of the interface
	netmask    	TEXT, 			-- Netmask of the interface
	is_dhcp		INTEGER NOT NULL DEFAULT 0 CHECK (is_dhcp IN (0, 1)),	-- Whether DHCP is enabled
	is_pxe		INTEGER NOT NULL DEFAULT 0 CHECK (is_pxe IN (0, 1)),	-- Whether PXE is enabled
	is_management	INTEGER NOT NULL DEFAULT 0 CHECK (is_management IN (0, 1)),	-- Whether this is a management interface
	type		TEXT,			-- Type of the interface (e.g., physical, virtual)
	bond_mode	TEXT,			-- Bond mode for bonded interfaces
	master		INTEGER,		-- Foreign key to another interface (if bonded)
	network		INTEGER,		-- Foreign key to the networks table
	host		INTEGER NOT NULL,	-- Foreign key to the hosts table
	UNIQUE (name, host),			-- Combination of name and host must be unique
	FOREIGN KEY (master) REFERENCES host_interfaces(id),
	FOREIGN KEY (network) REFERENCES networks(id),
	FOREIGN KEY (host) REFERENCES hosts(id) ON DELETE CASCADE
);

CREATE TABLE software (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the software
	type		TEXT NOT NULL,		-- Type of the software
	url		TEXT,			-- URL for downloading the software
	key_url		TEXT,			-- URL for downloading the software key
	architecture	TEXT,			-- Architecture of the software
	stack		INTEGER NOT NULL,	-- Foreign key to the stacks table
	FOREIGN KEY (stack) REFERENCES stacks(id) ON DELETE CASCADE
);

CREATE TABLE switches (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the switch
	name		TEXT NOT NULL,		-- Name of the switch
	model		INTEGER,		-- Foreign key to the models table
	environment	INTEGER,		-- Foreign key to the environments table
	rack		INTEGER,		-- Foreign key to the racks table
	rank		INTEGER,		-- Rank of the switch in the rack
	zone		INTEGER,		-- Foreign key to the zones table
	UNIQUE (name, zone),			-- Combination of name and zone must be unique
	FOREIGN KEY (model) REFERENCES models(id),
	FOREIGN KEY (environment) REFERENCES environments(id),
	FOREIGN KEY (rack) REFERENCES racks(id),
	FOREIGN KEY (zone) REFERENCES zones(id) ON DELETE CASCADE
);

CREATE TABLE switch_interfaces (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the switch interface
	name		TEXT NOT NULL,		-- Name of the interface
	switch		INTEGER NOT NULL,	-- Foreign key to the switches table
	UNIQUE (name, switch),			-- Combination of name and switch must be unique
	FOREIGN KEY (switch) REFERENCES switches(id) ON DELETE CASCADE
);




CREATE TABLE bmcs (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the BMC
	name		TEXT,			-- Name of the BMC
	host		INTEGER NOT NULL,	-- Foreign key to the hosts table
	model		INTEGER,		-- Foreign key to the models table
	environment	INTEGER,		-- Foreign key to the environments table
	UNIQUE (name, host),		-- Combination of name and host must be unique
	FOREIGN KEY (host) REFERENCES hosts(id) ON DELETE CASCADE,
	FOREIGN KEY (model) REFERENCES models(id),
	FOREIGN KEY (environment) REFERENCES environments(id)
);

CREATE TABLE bmc_interfaces (
	id		INTEGER PRIMARY KEY,	-- Unique identifier for the BMC interface
	name		TEXT NOT NULL,		-- Name of the interface
	ip		TEXT,			-- IP address of the interface
	mac		TEXT,			-- MAC address of the interface
	bmc		INTEGER NOT NULL,	-- Foreign key to the BMCs table
	network		INTEGER,		-- Foreign key to the networks table
	UNIQUE (name, bmc),			-- Combination of name and BMC must be unique
	FOREIGN KEY (bmc) REFERENCES bmcs(id) ON DELETE CASCADE,
	FOREIGN KEY (network) REFERENCES networks(id)
);


-- +goose Down

DROP TABLE bmc_interfaces;
DROP TABLE bmcs;

DROP TABLE host_interfaces;
DROP TABLE hosts;
DROP TABLE software;
DROP TABLE switch_interfaces;
DROP TABLE switches;

DROP TABLE appliances;
DROP TABLE attributes;
DROP TABLE clusters;
DROP TABLE environments;
DROP TABLE networks;
DROP TABLE racks;
DROP TABLE stacks;
DROP TABLE user_roles;

DROP TABLE entities;
DROP TABLE models;
DROP TABLE makes;
DROP TABLE oses;
DROP TABLE roles;
DROP TABLE users;
DROP TABLE zones;
