package system

import (
	"github.com/downflux/go-ecs/component"
	"github.com/downflux/go-ecs/id"
)

type S func(eid id.EID, cs map[id.CID]component.C) error
