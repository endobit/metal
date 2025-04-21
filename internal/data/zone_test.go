package data

import (
	"context"
	"testing"
)

func TestZones(t *testing.T) {
	s, close := setupDatabase(t)
	defer close(t)

	names := []string{"zone-a", "zone-b", "zone-c"}

	for _, zone := range names {
		_, err := s.CreateZone(context.Background(), zone)
		ok(t, err)
	}

	_, err := s.CreateZoneAttr(context.Background(), "attr-a", ZoneScope{Zone: "zone-a"})
	ok(t, err)
	err = s.UpdateZoneAttrValue(context.Background(), "value-a", ZoneAttrScope{Zone: "zone-a", Attr: "attr-a"})
	ok(t, err)

	t.Run("Read", func(t *testing.T) {
		zones, err := s.ReadZones(context.Background(), Filter{})
		ok(t, err)
		equals(t, len(zones), len(names))
		for i, zone := range zones {
			equals(t, names[i], zone.Name)
		}
	})

	t.Run("ReadAttrs", func(t *testing.T) {
		attrs, err := s.ReadZoneAttrs(context.Background(), Filter{}, ZoneScope{Zone: "zone-a"})
		ok(t, err)
		equals(t, len(attrs), 1)
		equals(t, attrs[0].Name, "attr-a")
		equals(t, attrs[0].Value, "value-a")
	})

}
