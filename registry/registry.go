package registry

import (
	"fmt"

	"github.com/downflux/go-ecs/component"
	"github.com/downflux/go-ecs/entity"
	"github.com/downflux/go-ecs/id"
	"github.com/downflux/go-ecs/system"
	"github.com/downflux/go-pq/pq"
)

type R struct {
	entities   map[id.EID]*entity.E
	components map[id.CID]map[id.EID]component.C
}

// Each allows the user to query and queue component mutations.
func (r *R) Each(components []id.CID, s system.S, pool int) {
	if len(components) == 0 {
		return
	}

	// Loop over the component with the least amount of elements, then look
	// up each element for that component map.
	q := pq.New[id.CID](len(components), pq.PMin)
	for _, cid := range components {
		q.Push(cid, float64(len(r.components[cid])))
	}

	candidates := make(map[id.EID]map[id.CID]component.C, 16)
	cid, _ := q.Pop()
	for eid, c := range r.components[cid] {
		candidates[eid] = make(map[id.CID]component.C, len(components))
		candidates[eid][cid] = c
	}

	if !q.Empty() {
		for cid, _ := q.Pop(); !q.Empty(); cid, _ = q.Pop() {
			// Entity candidates must contain all components.
			for eid, cs := range candidates {
				if c, ok := r.components[cid][eid]; !ok {
					delete(candidates, eid)
				} else {
					cs[cid] = c
				}
			}
		}
	}

	// TODO(minkezhang): Use pool size.
	for _, c := range candidates {
		s(c)
	}
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
