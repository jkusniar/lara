package api

import (
	"github.com/jkusniar/lara/msg"
)

func RegisterHandlers() *msg.HandlerRegistry {
	var registry = make(msg.HandlerRegistry)

	registry.RegisterHandler("CreateClient", CreateClient)

	return &registry
}
