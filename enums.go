package metal

import (
	"iter"
	"maps"
	"strings"

	pb "endobit.io/metal/gen/go/proto/metal/v1"
)

var (
	LongArchitecture   = longMapCreate("ARCHITECTURE", maps.Keys(pb.Architecture_value))
	ShortArchitecture  = shortMapCreate("ARCHITECTURE", maps.Keys(pb.Architecture_value))
	LongBondMode       = longMapCreate("BOND_MODE", maps.Keys(pb.BondMode_value))
	ShortBondMode      = shortMapCreate("BOND_MODE", maps.Keys(pb.BondMode_value))
	LongHostType       = longMapCreate("HOST_TYPE", maps.Keys(pb.HostType_value))
	ShortHostType      = shortMapCreate("HOST_TYPE", maps.Keys(pb.HostType_value))
	LongInterfaceType  = longMapCreate("INTERFACE_TYPE", maps.Keys(pb.InterfaceType_value))
	ShortInterfaceType = shortMapCreate("INTERFACE_TYPE", maps.Keys(pb.InterfaceType_value))
	LongOSFlavor       = longMapCreate("OS_FLAVOR", maps.Keys(pb.OSFlavor_value))
	ShortOSFlavor      = shortMapCreate("OS_FLAVOR", maps.Keys(pb.OSFlavor_value))
	LongOSName         = longMapCreate("OSNAME", maps.Keys(pb.OSName_value))
	ShortOSName        = shortMapCreate("OSNAME", maps.Keys(pb.OSName_value))
	LongSoftwareType   = longMapCreate("SOFTWARE_TYPE", maps.Keys(pb.SoftwareType_value))
	ShortSoftwareType  = shortMapCreate("SOFTWARE_TYPE", maps.Keys(pb.SoftwareType_value))
)

func shortName(prefix, name string) string {
	return strings.ToLower(strings.TrimPrefix(name, prefix+"_"))
}

func longMapCreate(prefix string, longNames iter.Seq[string]) map[string]string {
	m := make(map[string]string)

	for name := range longNames {
		m[shortName(prefix, name)] = name
	}

	return m
}

func shortMapCreate(prefix string, shortNames iter.Seq[string]) map[string]string {
	m := make(map[string]string)

	for name := range shortNames {
		m[name] = shortName(prefix, name)
	}

	return m
}
