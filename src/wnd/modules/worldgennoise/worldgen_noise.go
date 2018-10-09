package worldgennoise

import (
	"wnd/api"
	"wnd/utils/log"

	"wnd/modules"
	"wnd/modules/worldgennoise/layers"
)

type WorldGeneratorLayer interface {
	ID() string
	//Priority() int8
	Fill(seed int64, c *api.Chunk, size uint32) error
}

func New() *worldGenerator {
	return &worldGenerator{layers: make([]WorldGeneratorLayer, 0, 2)}
}

type worldGenerator struct {
	Universe      *api.Universe         `inject:""`
	BlockRegistry modules.BlockRegistry `inject:""`

	layers []WorldGeneratorLayer
}

func (this *worldGenerator) GenerateChunk(coords api.WorldCoords) *api.Chunk {
	c := &api.Chunk{Coords: coords}
	
	for _, l := range this.layers {
		log.Debugf("Result of %v: %v", l.ID(), l.Fill(this.Universe.Seed, c, this.Universe.Size * api.ChunkSideSize))
	}

	return c
}

func (this *worldGenerator) ID() string {
	return "worldGenerator"
}

func (this *worldGenerator) Start() error {
	return nil
}
func (this *worldGenerator) Stop() {
}
func (this *worldGenerator) Init() error {
	this.layers = append(this.layers, layers.NewBiomeLayer(), layers.NewTerrainLayer(this.Universe, this.BlockRegistry.All()))
	return nil
}
func (this *worldGenerator) Destroy() {

}
func (this *worldGenerator) Configuration() []api.TypeKeyValue {
	return nil
}
func (this *worldGenerator) SetConfiguration(values ...api.TypeKeyValue) error {
	return nil
}
