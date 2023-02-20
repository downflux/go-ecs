package registry

import (
	"testing"

	"github.com/downflux/go-ecs/component"
	"github.com/downflux/go-ecs/id"
)

func TestGet(t *testing.T) {
	r := New()
	for i := 0; i < 1000; i++ {
		r.Insert(map[id.CID]component.C{
			id.CID(i): i,
		})
	}
	for i := 0; i < 1000; i++ {
		eid := id.EID(i + 1)
		cid := id.CID(i)
		want := i
		if got := r.Get(eid, cid).(int); got != want {
			t.Errorf("Get() = %v, want = %v\n", got, want)

		}
	}
}

func TestEach(t *testing.T) {
	r := New()

	// User defines the static component IDs in a custom lookup table.
	const cEven = 0
	const cTriplet = 100

	for i := 0; i < 1000; i++ {
		r.Insert(map[id.CID]component.C{
			id.CID(i % 2):              i,
			id.CID(cTriplet + (i % 3)): i,
		})
	}

	t.Run("Evens", func(t *testing.T) {
		got := 0
		want := 249500
		r.Each([]id.CID{cEven}, func(eid id.EID, cs map[id.CID]component.C) error {
			// Entities 0, 2, ... 1000 are included in this sum.
			got += cs[cEven].(int)
			return nil
		}, 1)

		if got != want {
			t.Errorf("Each() = %v, want = %v\n", got, want)
		}
	})
	t.Run("MultiField", func(t *testing.T) {
		got := 0
		want := 2 * 83166
		// Filter to include only multiples of 6.
		r.Each([]id.CID{cEven, cTriplet}, func(eid id.EID, cs map[id.CID]component.C) error {
			got += cs[cEven].(int)
			got += cs[cTriplet].(int)
			return nil
		}, 1)

		if got != want {
			t.Errorf("Each() = %v, want = %v\n", got, want)
		}
	})
}
