package modules

import (
	"wnd/api"
)

type Broadcast interface {
	api.OutputModule
	
	SendTo(dst interface{}, e api.Event, reliable bool) error
}