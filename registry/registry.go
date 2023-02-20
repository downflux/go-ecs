package registry

import (
	"fmt"

	"github.com/downflux/go-ecs/component"
	"github.com/downflux/go-ecs/entity"
	"github.com/downflux/go-ecs/id"
	"github.com/downflux/go-ecs/system"
	// "github.com/downflux/go-pq/pq"
)

type R struct {
	entities   map[id.EID]*entity.E
	components map[id.CID]map[id.EID]component.C
}

// Each allows the user to query and queue component mutations.
func (r *R) Each(components []id.CID, s system.S) {
	// Loop over the component with the least amount of elements, then look
	// up each element for that component map; if CID matches AND does not
	// any in exclude...
}

func (r *R) Insert(cs map[id.CID]component.C) id.EID {
	eid := id.EID(len(r.entities) + 1)

	components := make([]id.CID, 0, len(cs))
	for cid, c := range cs {
		components = append(components, cid)
		if _, ok := r.components[cid]; !ok {
			r.components[cid] = make(map[id.EID]component.C, 16)
		}
		r.components[cid][eid] = c
	}

	e := entity.New(eid, components)
	r.entities[eid] = e

	return eid
}

func (r *R) Remove(eid id.EID) {
	e, ok := r.entities[eid]
	if !ok {
		panic(fmt.Sprintf("cannot remove non-existent entity %v", eid))
	}

	for _, cid := range e.Components() {
		delete(r.components[cid], eid)
	}
	delete(r.entities, eid)
}

func (r *R) Get(eid id.EID, cid id.CID) component.C {
	es, ok := r.components[cid]
	if !ok {
		panic(fmt.Sprintf("cannot get non-existent component %v", cid))
	}
	if c, ok := es[eid]; !ok {
		panic(fmt.Sprintf("cannot get non-existent entities %v", eid))
	} else {
		return c
	}
}
