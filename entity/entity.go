package entity

import (
	"github.com/downflux/go-ecs/id"
)

type E struct {
	id         id.EID
	components []id.CID
}

func New(eid id.EID, components []id.CID) *E {
	e := &E{
		id:         eid,
		components: components,
	}
	return e
}

func (e *E) ID() id.EID           { return e.id }
func (e *E) Components() []id.CID { return e.components }
