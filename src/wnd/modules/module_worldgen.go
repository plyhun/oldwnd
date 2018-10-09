package modules

import (
	"wnd/api"
)

type WorldGenerator interface {
	api.InittableModule
	api.ConfigurableModule
	
	GenerateChunk(coords api.WorldCoords) *api.Chunk
}