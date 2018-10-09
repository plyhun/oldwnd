package modules

import (
	"wnd/api"
)

type Storage interface {
	api.GameModule
	api.PrioritizedModule
	
	NewUniverse(u *api.Universe) error
	Universe(universeName string) *api.Universe
	DeleteUniverse(universeName string)
	
	SaveChunk(universeName string, chunk *api.Chunk)
	Chunk(universeName string, coords api.WorldCoords) (*api.Chunk,error)
	
	SaveEntity(universeName string, entity *api.Entity)
	Entity(universeName, entityId string) (*api.Entity,error)
}