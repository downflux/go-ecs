package entity

import (
	"github.com/downflux/go-ecs/id"
)

type E struct {
	eid id.EID
	cid id.CID
}

func New(eid id.EID, cid id.CID) *E {
	return &E{
		eid: eid,
		cid: cid,
	}
}

func (e *E) EID() id.EID { return e.eid }
func (e *E) CID() id.CID { return e.cid }
