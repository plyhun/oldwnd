package modules

import (
	"wnd/api"
	"wnd/api/events"
)

type UniverseProcessor interface {
	api.ProcessModule

	Chunks(coords []api.WorldCoords) []*api.Chunk
	Blocks(coords []api.WorldCoords) []*api.Block
	World(entityId string) *api.World
	Entity(entityId string) *api.Entity
	AddEntity(entity *api.Entity)
	MoveEntity(event *events.EntityMove) *events.EntityPosition
}
