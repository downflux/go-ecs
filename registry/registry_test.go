package registry

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/downflux/go-ecs/component"
	"github.com/downflux/go-ecs/id"
)

func BenchmarkEach(b *testing.B) {
	const (
		cBase   = 0
		cDense  = 1
		cSparse = 2
	)

	type config struct {
		n  int
		cs int
	}

	configs := []config{
		{
			n:  1e4,
			cs: 0,
		},
		{
			n:  1e5,
			cs: 0,
		},
		{
			n:  1e6,
			cs: 0,
		},
		{
			n:  1e4,
			cs: 100,
		},
		{
			n:  1e5,
			cs: 100,
		},
		{
			n:  1e6,
			cs: 100,
		},
	}

	for _, c := range configs {
		r := New()
		for i := 0; i < c.n; i++ {
			cs := map[id.CID]component.C{
				cBase: true,
			}
			if rand.Float64() > 0.3 {
				cs[cDense] = true
			}
			if rand.Float64() > 0.75 {
				cs[cSparse] = true
			}

			for j := 0; j < c.cs; j++ {
				cs[id.CID(j)+2] = true
			}

			r.Insert(cs)
		}
		b.Run(fmt.Sprintf("Dense/N=%v/C=%v", c.n, c.cs+3), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r.Each([]id.CID{cBase, cDense}, func(eid id.EID, cs map[id.CID]component.C) error { return nil }, 1)
			}
		})
		b.Run(fmt.Sprintf("Sparse/N=%v/C=%v", c.n, c.cs+3), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r.Each([]id.CID{cBase, cSparse}, func(eid id.EID, cs map[id.CID]component.C) error { return nil }, 1)
			}
		})
	}
}

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
