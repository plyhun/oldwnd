package events

import (
	"wnd/api"
)

type General struct {
	EventSource     string
	EventMetadata   interface{}
	EventID         string
	EventTime       uint64
	EventUniverseID string
}

type Outputable struct {
	EventTarget string
}

func (this *General) ID() string {
	return this.EventID
}

func (this *General) Time() uint64 {
	return this.EventTime
}

func (this *General) UniverseID() string {
	return this.EventUniverseID
}

func (this *General) Metadata() interface{} {
	return this.EventMetadata
}

func (this *General) SetMetadata(src interface{}) {
	this.EventMetadata = src
}

func (this *General) Source() string {
	return this.EventSource
}

func (this *Outputable) Target() string {
	return this.EventTarget
}

type World struct {
	General
	Outputable

	World      *api.World
	Radius uint8
	NeedChunks bool
}

type Chunk struct {
	General
	Outputable

	NoBlocks bool
	Coords   api.WorldCoords
	Chunk    *api.Chunk
}

type Chunks struct {
	General
	Outputable

	NoBlocks, SplitOutput bool
	Coords   []api.WorldCoords
	Chunks    []*api.Chunk
}

type Blocks struct {
	General
	Outputable
	Blocks map[api.PackedWorldCoords]*api.Block
}

type EntityMove struct {
	General
	Outputable
	EntityID                       string
	DeltaTime                      int64
	HorizontalAngle, VerticalAngle float64
	Offset                         api.Coords
}

type EntityPosition struct {
	General
	Outputable
	api.Observer
	
	EntityID string
	IsPlayer bool
}

type EntityChanged struct {
	General
	Entity *api.Entity
}

type CustomBlocks struct {
	General
	Outputable

	Blocks []api.BlockDefinition
}
