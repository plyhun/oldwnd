package innerapi

import (
	"wnd/api"

	glm "github.com/go-gl/mathgl/mgl32"
)

type RenderDataID string

type RenderProgramID string

type RenderData interface {
	api.Outputable

	RenderDataID() RenderDataID

	Render(camera *GlCamera)
	Purge()
}

type FOVRenderData interface {
	RenderData

	IsVisible(camera *GlCamera, w *api.World) bool
}

type ProgrammableRenderData interface {
	RenderData

	Init(programs map[RenderProgramID]RenderProgram)
	RenderProgramIDs() []RenderProgramID
}

type RenderProgram struct {
	Program       uint32
	Texture       uint32
	NormalTexture uint32

	MVP         int32
	V           int32
	M           int32
	LightDir    int32
	SkyColor    int32
	GroundColor int32
	Shift int32

	Texs int32
	Norm int32

	Color    uint32
	Vertex   uint32
	Normal   uint32
	Position uint32
}

type GlCamera struct {
	*api.Observer

	M, V, P, MV, MVP glm.Mat4
}

func (this *GlCamera) ObserverData() *api.Observer {
	return this.Observer
}
