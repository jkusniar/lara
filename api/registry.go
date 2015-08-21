package api

import (
	"github.com/jkusniar/lara/msg"
)

func RegisterHandlers() *msg.Registry {
	var reg = make(msg.Registry)

	reg.Register("CreateClient", CreateClient)

	return &reg
}
