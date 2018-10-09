package graphics

import (
	"wnd/api"
	
	glm "github.com/go-gl/mathgl/mgl32"
)

type RemoveOutputable struct {
	OutputableID string
}

var (
	UD     = glm.Vec3{0, 1, 0}
	TARGET = "modules.Graphics"
)

func (this *RemoveOutputable) ID() string {
	return this.OutputableID
}

func (this *RemoveOutputable) Target() string {
	return TARGET
}

func (this *RemoveOutputable) IsConsumed() bool {
	return false
}

func (this *RemoveOutputable) Consume() {
}

func GetGraphicsOutputableID(src interface{}) string {
	switch t := src.(type) {
	case string:
		return t
	case *api.Chunk:
		return "chunk " + string(t.Coords.Pack())
	case *api.Block:
		return "block " + string(t.PackWithoutCoords())
	case *api.BlockInChunk:
		return "block " + string(t.PackWithoutCoords())
	}

	return ""
}

type Camera interface {
	ObserverData() *api.Observer
}
