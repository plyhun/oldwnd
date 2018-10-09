package modules

import (
	"wnd/api"
)

type Tattva interface {
	api.GameModule
	
	NewUniverse(u *api.Universe) error
	UniverseList() []*api.Universe
	Universe(universeId string) *api.Universe
	HasUniverse(universeId string, stored bool) bool
	RemoveUniverse(universeId string, fromStorage bool)
}