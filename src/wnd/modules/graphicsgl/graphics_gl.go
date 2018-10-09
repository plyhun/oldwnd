package graphicsgl

import (
	"wnd/api"
	"wnd/api/events"
	"wnd/base/graphics"
	"wnd/modules"
	"wnd/modules/graphicsgl/glutils"
	"wnd/modules/graphicsgl/innerapi"
	"wnd/modules/graphicsgl/renderers"
	"wnd/utils/log"

	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw3/v3.1/glfw"
	glm "github.com/go-gl/mathgl/mgl32"
)

const (
	_MAX_FPS = 100
)

type graphicsGl struct {
	window *glfw.Window
	camera *innerapi.GlCamera

	observerChannel          chan *api.Observer
	renderablesChannel       chan innerapi.RenderData
	renderablesDeleteChannel chan innerapi.RenderDataID

	renderables, visibleRenderables *glutils.OrderedMap

	Context       modules.GlobalContext `inject:""`
	BlockRegistry modules.BlockRegistry `inject:""`
	World         *api.World            `inject:""`

	time int64
}

func New() modules.Graphics {
	return &graphicsGl{
		renderables:              glutils.NewOrderedMap(),
		visibleRenderables:       glutils.NewOrderedMap(),
		observerChannel:          make(chan *api.Observer),
		renderablesChannel:       make(chan innerapi.RenderData),
		renderablesDeleteChannel: make(chan innerapi.RenderDataID),
		camera: &innerapi.GlCamera{
			//P: glm.Ortho(-100, 100,-100, 100,-100, 100),
			M: glm.Ident4(),
		},
	}
}

func (this *graphicsGl) ID() string {
	return "graphicsGl"
}

func (this *graphicsGl) Priority() int8 {
	return -128
}

func (this *graphicsGl) Start() (err error) {
	runtime.LockOSThread()

	if err = glfw.Init(); err != nil {
		panic("Failed to initialize GLFW")
	}

	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	pm := glfw.GetPrimaryMonitor()
	vm := pm.GetVideoMode()

	this.window, err = glfw.CreateWindow(2*vm.Width/3, 2*vm.Height/3, "wnd", nil, nil)
	if err != nil {
		glfw.Terminate()
		return err
	}

	this.window.MakeContextCurrent()

	if init := gl.Init(); init != nil {
		panic(init)
	}

	this.window.SetInputMode(glfw.StickyKeysMode, glfw.True)
	this.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	gl.ClearColor(renderers.SKY[0], renderers.SKY[1], renderers.SKY[2], 1.0)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	this.Context.Put("window", this.window)

	this.time = time.Now().UnixNano()

	return nil
}

func (this *graphicsGl) Init() error {
	this.camera.P = glm.Perspective(glm.DegToRad(60.0), 4.0/3.0, 0.1, float32(api.ChunkSideSize)*float32(this.World.VisibleRadius + 1)*2)
	return nil
}

func (this *graphicsGl) Stop() {
	if this.window != nil {
		this.window.Destroy()
	}

	runtime.UnlockOSThread()
}

func (this *graphicsGl) Destroy() {
	glfw.Terminate()
	this.window = nil
}

func (this *graphicsGl) Output(time uint64, renderables []api.Outputable) {
	for _, r := range renderables {
		if r.Target() == "" || r.Target() == graphics.TARGET {
			log.Debugf("%#v", renderables)
			this.addInternal(r)
		}
	}

	this.addExternal()
	this.draw()
}

func (this *graphicsGl) addExternal() {
	select {
	case mv := <-this.observerChannel:
		this.viewerPosition(mv)
	case rnd := <-this.renderablesChannel:
		this.renderData(rnd)
	case drnd := <-this.renderablesDeleteChannel:
		this.removeOutputable(drnd)
	default:
		//log.Tracef("No render update")
	}
}

func (this *graphicsGl) viewerPosition(o *api.Observer) {
	log.Tracef("%v", o)

	this.camera.Observer = o
	this.camera.V = glm.LookAtV(glm.Vec3{o.Position.X, o.Position.Y, o.Position.Z}, glm.Vec3{o.LooksAt.X, o.LooksAt.Y, o.LooksAt.Z}, graphics.UD)
	this.camera.MV = this.camera.V.Mul4(this.camera.M)
	this.camera.MVP = this.camera.P.Mul4(this.camera.MV)

	this.refreshVisibles()
}

func (this *graphicsGl) removeOutputable(drnd innerapi.RenderDataID) {
	log.Debugf("delete %v vr %v r %v", drnd, this.visibleRenderables.Delete(drnd) != nil, this.renderables.Delete(drnd) != nil)
}

func (this *graphicsGl) refreshVisibles() {
	this.visibleRenderables.Clear()

	for _, k := range this.renderables.Keys() {
		r, _ := this.renderables.Get(k)

		if fr, ok := r.(innerapi.FOVRenderData); ok {
			if fr.IsVisible(this.camera, this.World) {
				this.visibleRenderables.Put(fr)
			}
		} else {
			this.visibleRenderables.Put(r)
		}
	}

	//log.Warnf("%v of %v visible", len(this.visibleRenderables), len(this.renderables))
}

func (this *graphicsGl) addInternal(r api.Outputable) {
	log.Tracef("%#v", r)

	switch tr := r.(type) {
	case innerapi.RenderData:
		this.renderData(tr)
	case *events.EntityPosition:
		if tr.IsPlayer {
			this.viewerPosition(&tr.Observer)
		}
	case *graphics.RemoveOutputable:
		this.removeOutputable(innerapi.RenderDataID(tr.OutputableID))
	case *events.Chunk:
		if tr.Chunk != nil {
			go this.processChunk(tr.Chunk)
		}
	case *events.Chunks:
		go func(e *events.Chunks) {
			for _, c := range e.Chunks {
				c := this.World.Chunks.Data[c.Coords.Pack()]
				if c != nil {
					this.processChunk(c)
				}
			}
		}(tr)
	case *events.World:
		//go this.trace(tr.World)
	}
}

func (this *graphicsGl) processChunk(c *api.Chunk) {
	crd := renderers.CreateChunkRenderData(c, this.World)

	if crd != nil {
		//log.Warnf("%v <= %v", crd.ID(), c.Coords)

		this.AddToRender(crd)
	}
}

func (this *graphicsGl) renderData(rnd innerapi.RenderData) {
	log.Tracef("rnd <- %#v", rnd)

	if rr := this.renderables.Delete(rnd.RenderDataID()); rr != nil {
		rr.Purge()
	}
	this.visibleRenderables.Delete(rnd.RenderDataID())

	if p, ok := rnd.(innerapi.ProgrammableRenderData); ok {
		ids := p.RenderProgramIDs()

		prs := make(map[innerapi.RenderProgramID]innerapi.RenderProgram, len(ids))

		for _, id := range ids {
			prs[id] = renderers.GetProgramForBlockType(id, this.BlockRegistry)
		}

		p.Init(prs)
	}

	this.renderables.Put(rnd)

	fovOwner := this.World.Entities[this.World.Owner]

	if fovOwner != nil {
		if fr, ok := rnd.(innerapi.FOVRenderData); ok {
			if fr.IsVisible(this.camera, this.World) {
				this.visibleRenderables.Put(fr)
			}
		} else {
			this.visibleRenderables.Put(rnd)
		}
	}

	log.Debugf("%v of %v visible", this.visibleRenderables.Size(), this.renderables.Size())
}

func (this *graphicsGl) AddToRender(r api.Outputable) {
	log.Tracef("%#v", r)

	switch tr := r.(type) {
	case innerapi.RenderData:
		this.renderablesChannel <- tr
	case *events.EntityPosition:
		if tr.IsPlayer {
			this.observerChannel <- &tr.Observer
		}
	case *graphics.RemoveOutputable:
		this.RemoveFromRender(tr.OutputableID)
	case *events.Chunk:
		if tr.Chunk != nil {
			go this.processChunk(tr.Chunk)
		}
	case *events.Chunks:
		go func(e *events.Chunks) {
			for _, c := range e.Chunks {
				c := this.World.Chunks.Data[c.Coords.Pack()]
				if c != nil {
					this.processChunk(c)
				}
			}
		}(tr)
	case *events.World:
		//go this.trace(tr.World)
	}
}

func (this *graphicsGl) RemoveFromRender(rid string) {
	log.Tracef("%v", rid)
	this.renderablesDeleteChannel <- innerapi.RenderDataID(rid)
}

func (this *graphicsGl) draw() {
	end := time.Now().UnixNano()

	delta := end - this.time
	if delta == 0 {
		delta = 1
	} else if delta > int64(2*time.Second) {
		delta = int64(2 * time.Second)
	}

	fps := 1000000000 / delta

	if fps > _MAX_FPS || this.window == nil || this.window.ShouldClose() {
		return
	}

	this.time = end

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, k := range this.visibleRenderables.Keys() {
		v, _ := this.visibleRenderables.Get(k)
		log.Debugf("Rendering %s", k)
		v.Render(this.camera)
	}

	this.window.SwapBuffers()
	glfw.PollEvents()
	
	ww, wh := this.window.GetSize()
	
	var depth float32
	
	gl.ReadPixels(int32(ww) / 2, int32(wh) / 2, 1, 1, gl.DEPTH_COMPONENT, gl.FLOAT, gl.Ptr(&depth))
	
	wincoord := glm.Vec3{float32(ww) / 2, float32(wh) / 2, depth};
	objcoord,err := glm.UnProject(wincoord, this.camera.MV, this.camera.P, 0, 0, ww, wh);
	
	if err != nil {
		log.Errorf("%v", err)
		this.window.SetTitle(fmt.Sprintf("wnd: %v fps, %v of %v visible", fps, this.visibleRenderables.Size(), this.renderables.Size()))
	} else {
		this.window.SetTitle(fmt.Sprintf("wnd: %v fps, %v of %v visible (center at %v/%v/%v)", fps, this.visibleRenderables.Size(), this.renderables.Size(), int(objcoord[0]), int(objcoord[1]), int(objcoord[2])))
	}
	
	//this.window.SetTitle(fmt.Sprintf("wnd: %v fps, %v of %v visible", fps, this.visibleRenderables.Size(), this.renderables.Size()))
	
	//runtime.Gosched()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
}
