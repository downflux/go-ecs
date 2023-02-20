package registry

import (
	"github.com/downflux/go-ecs/component"
	"github.com/downflux/go-ecs/entity"
	"github.com/downflux/go-ecs/id"
	"github.com/downflux/go-ecs/system"
)

type R struct {
	entities   map[id.EID]*entity.E
	components map[id.CID]map[id.EID]component.C
}

// Each allows the user to query and queue component mutations. Mutations must
// not be committed (i.e. c.Commit) here.
func (r *R) Each(cid id.CID, exclude id.CID, s system.S) {
	// Loop over the component with the least amount of elements, then look
	// up each element for that component map; if CID matches AND does not
	// any in exclude...
}

func (r *R) Insert(cs map[id.CID]component.C) id.EID {
	eid := id.EID(len(r.entities) + 1)
	cid := id.CIDNone
	for k, v := range cs {
		cid |= k
		if _, ok := r.components[k]; !ok {
			r.components[k] = make(map[id.EID]component.C, 1024)
		}
		r.components[k][eid] = v
	}

	e := entity.New(eid, cid)
	r.entities[eid] = e

	return eid
}

func (r *R) Remove(es ...id.EID) {
	for _, eid := range es {
		e := r.entities[eid]
		cid := e.CID()

		for i := 1; i < 64; i++ {
			if mask := id.CID(1 << i); cid&mask != id.CIDNone {
				delete(r.components[cid], eid)
			}
		}
		delete(r.entities, eid)
	}
}

func (r *R) Get(eid id.EID, cid id.CID) component.C { return r.components[cid][eid] }
