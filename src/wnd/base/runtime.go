package base

import (
	"wnd/api"
	
	"reflect"
)

type Runtime interface {
	ID() string
	
	IsClient() bool
	IsServer() bool
	
	AddModule(m api.GameModule) error
	RemoveModule(id string)
	GetByType(t reflect.Type) []api.GameModule
	Start() error
	Stop() error 
}