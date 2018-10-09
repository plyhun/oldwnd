package modules

import (
	"wnd/api"
)

type EntityStorage interface {
	api.GameModule
	
	Add(e *api.Entity) error
	Remove(id string) error
	GetByID(id string) (*api.Entity,error)
	GetByCoords(coords api.WorldCoords, radius uint32) ([]*api.Entity,error)
	GetByControl(control string) ([]*api.Entity,error)
}