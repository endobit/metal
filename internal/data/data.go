// Package data provides the data store for the metal service. It abstracts the
// database layer and provides the data types and methods to interact with the
// database.
//
// New structs are used to copy out data from all database queries. This is done
// to hide all details of the database. For example: in SQLite booleans are
// represented are int64 values, and at this layer is bools.
package data

import (
	"database/sql"
	"embed"
	"log/slog"

	"endobit.io/metal/internal/data/db"
)

//go:embed migrations/*.sql
var Migrations embed.FS

type (
	User = db.User
	Role = db.Role

	// GlobalAttrScope qualifies queries that operate on a global attribute.
	GlobalAttrScope struct {
		Attr string
	}

	// MakeScope qualifies queries that operate on a make.
	MakeScope struct {
		Make string
	}

	// ModelScope qualifies queries that operate on a model.
	ModelScope struct {
		Make  string
		Model string
	}

	// ModelAttrScope qualifies queries that operate on a model attribute.
	ModelAttrScope struct {
		Make  string
		Model string
		Attr  string
	}

	// ZoneScope qualifies queries that operate on a zone.
	ZoneScope struct {
		Zone string
	}

	// ZoneAttrScope qualifies queries that operate on a zone attribute.
	ZoneAttrScope struct {
		Zone string
		Attr string
	}

	// NetworkScope qualifies queries that operate on a network.
	NetworkScope struct {
		Zone    string
		Network string
	}

	// RackScope qualifies queries that operate on a rack.
	RackScope struct {
		Zone string
		Rack string
	}

	// RackAttrScope qualifies queries that operate on a rack attribute.
	RackAttrScope struct {
		Zone string
		Rack string
		Attr string
	}

	// ApplianceScope qualifies queries that operate on an appliance.
	ApplianceScope struct {
		Zone      string
		Appliance string
	}

	// ApplianceAttrScope qualifies queries that operate on an appliance attribute.
	ApplianceAttrScope struct {
		Zone      string
		Appliance string
		Attr      string
	}

	// EnvironmentScope qualifies queries that operate on an environment.
	EnvironmentScope struct {
		Zone        string
		Environment string
	}

	// EnvironmentAttrScope qualifies queries that operate on an environment attribute.
	EnvironmentAttrScope struct {
		Zone        string
		Environment string
		Attr        string
	}

	// ClusterScope qualifies queries that operate on a cluster.
	ClusterScope struct {
		Zone    string
		Cluster string
	}

	// ClusterAttrScope qualifies queries that operate on a cluster attribute.
	ClusterAttrScope struct {
		Zone    string
		Cluster string
		Attr    string
	}

	// HostScope qualifies queries that operate on a host.
	HostScope struct {
		Zone    string
		Cluster string
		Host    string
	}

	// HostAttrScope qualifies queries that operate on a host attribute.
	HostAttrScope struct {
		Zone    string
		Cluster string
		Host    string
		Attr    string
	}

	// HostInterfaceScope qualifies queries that operate on a host interface.
	HostInterfaceScope struct {
		Zone      string
		Cluster   string
		Host      string
		Interface string
	}

	// Make represents a hardware manufacturer.
	Make struct {
		Name string
		id   int64
	}

	// Model represents a hardware model.
	Model struct {
		MakeScope
		Name         string
		Architecture string
		id           int64
	}

	// Zone represents a data center zone.
	Zone struct {
		Name     string
		TimeZone string
		id       int64
	}

	// Appliance represents a software profile. All Hosts have an Appliance type
	// that may include Appliance level attributes.
	Appliance struct {
		ZoneScope
		Name string
		id   int64
	}

	// Environment represents a software environment.
	Environment struct {
		ZoneScope
		Name string
		id   int64
	}

	// Rack represents a physical rack in a data center zone.
	Rack struct {
		ZoneScope
		Name string
		id   int64
	}

	// Network represents a network in a data center zone.
	Network struct {
		ZoneScope
		Name    string
		Address string
		Gateway string
		IsPXE   bool
		MTU     uint32
		id      int64
	}

	// Cluster represents a clustered set of Hosts.
	Cluster struct {
		ZoneScope
		Name string
		id   int64
	}

	// Host represents a physical or virtual machine.
	Host struct {
		ClusterScope
		Name        string
		Make        string
		Model       string
		Environment string
		Appliance   string
		Location    string
		Rack        string
		Type        string
		rank        *uint32
		slot        *uint32
		id          int64
	}

	// Attr represents a string key value pair.
	Attr struct {
		Name        string
		Value       string
		IsProtected bool
		id          int64
	}

	// GlobalAttr represents a global attribute.
	GlobalAttr struct {
		Attr
	}

	// ModelAttr is a fully qualified model attribute.
	ModelAttr struct {
		ModelScope
		Attr
	}

	// ZoneAttr is a fully qualified zone attribute.
	ZoneAttr struct {
		ZoneScope
		Attr
	}

	// ApplianceAttr is a fully qualified appliance attribute.
	ApplianceAttr struct {
		ApplianceScope
		Attr
	}

	// EnvironmentAttr is a fully qualified environment attribute.
	EnvironmentAttr struct {
		EnvironmentScope
		Attr
	}

	// RackAttr is a fully qualified rack attribute.
	RackAttr struct {
		RackScope
		Attr
	}

	// ClusterAttr is a fully qualified cluster attribute.
	ClusterAttr struct {
		ClusterScope
		Attr
	}

	// HostAttr is a fully qualified host attribute.
	HostAttr struct {
		HostScope
		Attr
	}

	// HostInterface represents a network interface on a Host.
	HostInterface struct {
		HostScope
		Name            string
		IP              string
		MAC             string
		Netmask         string
		IsDHCP          bool
		IsPXE           bool
		IsManagement    bool
		Type            string
		BondMode        string
		MasterInterface string
		Network         string
		id              int64
	}

	// Store provides the data store for the metal service.
	Store struct {
		db *db.Queries
	}
)

// Rank returns the rank of the host if it is set. If the rank is not set, it
// returns 0, false.
func (h *Host) Rank() (rank uint32, ok bool) {
	if h.rank == nil {
		return 0, false
	}

	return *h.rank, true
}

// Slot returns the slot of the host if it is set. If the slot is not set, it
// returns 0, false.
func (h *Host) Slot() (slot uint32, ok bool) {
	if h.slot == nil {
		return 0, false
	}

	return *h.slot, true
}

// NewStore returns an initialized data store.
func NewStore(logger *slog.Logger, database *sql.DB) *Store {
	return &Store{
		db: db.New(loggingDB{
			DB:     database,
			Logger: logger,
			Level:  slog.LevelDebug,
		}),
	}
}

type FilterType int

const (
	NoFilter FilterType = iota
	NameFilter
	GlobFilter
)

// Filter represents a filter for a query. Filters can names or shell glob
// patterns.
type Filter struct {
	Type    FilterType
	Pattern string
}

// FilterByName returns a filter that matches an exact name.
func FilterByName(name string) Filter {
	return Filter{Type: NameFilter, Pattern: name}
}

// FilterByGlob returns a filter that matches a shell glob pattern.
func FilterByGlob(glob string) Filter {
	return Filter{Type: GlobFilter, Pattern: glob}
}

// Ptr returns a pointer to the value.
func Ptr[T any](v T) *T {
	return &v
}

// Optional return a pointer to the value if it is not the zero value. Otherwise
// it returns nil.
func Optional[T comparable](value T) *T {
	var zero T

	if value == zero {
		return nil
	}
	return &value
}

// value returns the value of a pointer if it is not nil. Otherwise it returns the
// zero value of the type.
func value[T comparable](value *T) T {
	var zero T

	if value == nil {
		return zero
	}
	return *value
}

// boolean returns value as a sqlite boolean (int64).
func boolean(value bool) int64 {
	if value {
		return 1
	}
	return 0
}
