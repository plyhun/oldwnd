package modules

import (
"wnd/api"
)

type Incoming interface {
	api.EventModule
	
	Parse(from,what interface{}) api.Event
}