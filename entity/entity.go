package entity

import (
	"github.com/downflux/go-ecs/id"
)

type E struct {
	eid        id.EID
	components []id.CID
}

func New(eid id.EID, components []id.CID) *E {
	e := &E{
		eid:        eid,
		components: components,
	}
	return e
}

func (e *E) EID() id.EID          { return e.eid }
func (e *E) Components() []id.CID { return e.components }
