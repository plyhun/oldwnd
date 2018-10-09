package renderers

import (

)

import (
	"wnd/api"
	"wnd/modules"
	"wnd/modules/graphicsgl/innerapi"
	"wnd/modules/graphicsgl/glutils"
	"wnd/utils/log"
	
	"strings"
	"strconv"
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/jragonmiris/mathgl"
)

var (
	FOGCOLOR = mathgl.Vec4f{0.95, 0.98, 1.0, 1.0} //rgba
)

type TypeSizeSlope string

var verticesSouth = [18]float32{
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.0, 0.0,
}

var normalsSouth = [18]float32{
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
}

var verticesNorth = [18]float32{
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.0, 0.9999,
}

var normalsNorth = [18]float32{
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
}

var verticesEast = [18]float32{
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
}

var normalsEast = [18]float32{
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
}	

var verticesWest = [18]float32{
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
}

var normalsWest = [18]float32{
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
}	

var verticesTop = [18]float32{
	0.0, 0.9999, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.9999, 0.9999,
}

var normalsTop = [18]float32{
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
}	

var verticesBottom = [18]float32{
	0.0, 0.0, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.0, 0.9999,
	0.0, 0.0, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.0, 0.0,
}

var normalsBottom = [18]float32{
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
}	

var verticesTopSouth = [18]float32{
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.0,
}

var normalsTopSouth = [18]float32{
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
}	

var verticesTopNorth = [18]float32{
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.9999,
}

var normalsTopNorth = [18]float32{
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
}	

var verticesTopEast = [18]float32{
	0.0, 0.0, 0.0, 
	0.9999, 0.9999, 0.0, 
	0.9999, 0.9999, 0.9999, 
	0.0, 0.0, 0.0, 
	0.9999, 0.9999, 0.9999, 
	0.0, 0.0, 0.9999, 
}

var normalsTopEast = [18]float32{
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
}	

var verticesTopWest = [18]float32{
	0.9999, 0.0, 0.9999, 
	0.0, 0.9999, 0.9999, 
	0.0, 0.9999, 0.0, 
	0.9999, 0.0, 0.9999, 
	0.0, 0.9999, 0.0, 
	0.9999, 0.0, 0.0, 
}

var normalsTopWest = [18]float32{
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
}	

var verticesBottomSouth = [18]float32{
	0.9999, 0.0, 0.9999, 
	0.9999, 0.9999, 0.0, 
	0.0, 0.9999, 0.0, 
	0.9999, 0.0, 0.9999, 
	0.0, 0.9999, 0.0, 
	0.0, 0.0, 0.9999, 
}

var normalsBottomSouth = [18]float32{
	0.0, -1.0, -1.0,
	0.0, -1.0, -1.0,
	0.0, -1.0, -1.0,
	0.0, -1.0, -1.0,
	0.0, -1.0, -1.0,
	0.0, -1.0, -1.0,
}	

var verticesBottomNorth = [18]float32{
	0.0, 0.0, 0.0, 
	0.0, 0.9999, 0.9999, 
	0.9999, 0.9999, 0.9999, 
	0.0, 0.0, 0.0, 
	0.9999, 0.9999, 0.9999, 
	0.9999, 0.0, 0.0, 
}

var normalsBottomNorth = [18]float32{
	0.0, -1.0, 1.0,
	0.0, -1.0, 1.0,
	0.0, -1.0, 1.0,
	0.0, -1.0, 1.0,
	0.0, -1.0, 1.0,
	0.0, -1.0, 1.0,
}	

var verticesBottomEast = [18]float32{
	0.9999, 0.0, 0.0, 
	0.0, 0.9999, 0.0, 
	0.0, 0.9999, 0.9999, 
	0.9999, 0.0, 0.0, 
	0.0, 0.9999, 0.9999, 
	0.9999, 0.0, 0.9999, 
}

var normalsBottomEast = [18]float32{
	-1.0, -1.0, 0.0,
	-1.0, -1.0, 0.0,
	-1.0, -1.0, 0.0,
	-1.0, -1.0, 0.0,
	-1.0, -1.0, 0.0,
	-1.0, -1.0, 0.0,
}	

var verticesBottomWest = [18]float32{
	0.0, 0.0, 0.9999, 
	0.9999, 0.9999, 0.9999, 
	0.9999, 0.9999, 0.0, 
	0.0, 0.0, 0.9999, 
	0.9999, 0.9999, 0.0, 
	0.0, 0.0, 0.0, 
}

var normalsBottomWest = [18]float32{
	1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
}	

var verticesSouthEast = [18]float32{
	0.9999, 0.0, 0.0, 
	0.9999, 0.9999, 0.0, 
	0.0, 0.9999, 0.9999, 
	0.9999, 0.0, 0.0, 
	0.0, 0.9999, 0.9999, 
	0.0, 0.0, 0.9999, 
}

var normalsSouthEast = [18]float32{
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
}	

var verticesSouthWest = [18]float32{
	0.9999, 0.0, 0.9999, 
	0.9999, 0.9999, 0.9999, 
	0.0, 0.9999, 0.0, 
	0.9999, 0.0, 0.9999, 
	0.0, 0.9999, 0.0, 
	0.0, 0.0, 0.0, 
}

var normalsSouthWest = [18]float32{
	1.0, 0.0, -1.0,
	1.0, 0.0, -1.0,
	1.0, 0.0, -1.0,
	1.0, 0.0, -1.0,
	1.0, 0.0, -1.0,
	1.0, 0.0, -1.0,
}	

var verticesNorthEast = [18]float32{
	0.0, 0.0, 0.0, 
	0.0, 0.9999, 0.0, 
	0.9999, 0.9999, 0.9999, 
	0.0, 0.0, 0.0, 
	0.9999, 0.9999, 0.9999, 
	0.9999, 0.0, 0.9999, 
}

var normalsNorthEast = [18]float32{
	-1.0, 0.0, 1.0,
	-1.0, 0.0, 1.0,
	-1.0, 0.0, 1.0,
	-1.0, 0.0, 1.0,
	-1.0, 0.0, 1.0,
	-1.0, 0.0, 1.0,
}	

var verticesNorthWest = [18]float32{
	0.0, 0.0, 0.9999, 
	0.0, 0.9999, 0.9999, 
	0.9999, 0.9999, 0.0, 
	0.0, 0.0, 0.9999, 
	0.9999, 0.9999, 0.0, 
	0.9999, 0.0, 0.0, 
}

var normalsNorthWest = [18]float32{
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
}	

var verticesTopSouthEast = [9]float32{
	0.0, 0.0, 0.0, 
	0.9999, 0.9999, 0.0, 
	0.0, 0.9999, 0.9999, 
}

var verticesTopSouthWest = [9]float32{
	0.9999, 0.0, 0.0, 
	0.9999, 0.9999, 0.9999, 
	0.0, 0.9999, 0.0, 
}

var verticesTopNorthEast = [9]float32{
	0.0, 0.0, 0.9999, 
	0.0, 0.9999, 0.0, 
	0.9999, 0.9999, 0.9999, 
}

var verticesTopNorthWest = [9]float32{
	0.9999, 0.0, 0.9999, 
	0.0, 0.9999, 0.9999, 
	0.9999, 0.9999, 0.0, 
}

var verticesBottomSouthEast = [9]float32{
	0.9999, 0.0, 0.0, 
	0.0, 0.9999, 0.0, 
	0.0, 0.0, 0.9999, 
}

var verticesBottomSouthWest = [9]float32{
	0.9999, 0.0, 0.9999, 
	0.9999, 0.9999, 0.0, 
	0.0, 0.0, 0.0, 
}

var verticesBottomNorthEast = [9]float32{
	0.0, 0.0, 0.0, 
	0.0, 0.9999, 0.9999, 
	0.9999, 0.0, 0.9999, 
}

var verticesBottomNorthWest = [9]float32{
	0.0, 0.0, 0.9999, 
	0.9999, 0.9999, 0.9999, 
	0.9999, 0.0, 0.0, 
}

var uvsNorth = [12]float32{
	0.6667, 0.4999,
	0.6667, 0.0,
	1.0, 0.0,
	0.6667, 0.4999,
	1.0, 0.0,
	1.0, 0.4999,
}

var uvsSouth = [12]float32{
	0.0, 0.4999,
	0.0, 0.0,
	0.3332, 0.0,
	0.0, 0.4999,
	0.3332, 0.0,
	0.3332, 0.4999,
}

var uvsWest = [12]float32{
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
}

var uvsEast = [12]float32{
	0.0, 1.0,
	0.0, 0.5001,
	0.3332, 0.5001,
	0.0, 1.0,
	0.3332, 0.5001,
	0.3332, 1.0,
}

var uvsTop = [12]float32{
	0.3334, 0.4999,
	0.3334, 0.0,
	0.6665, 0.0,
	0.3334, 0.4999,
	0.6665, 0.0,
	0.6665, 0.4999,
}

var uvsBottom = [12]float32{
	0.3334, 1.0,
	0.3334, 0.5001,
	0.6665, 0.5001,
	0.3334, 1.0,
	0.6665, 0.5001,
	0.6665, 1.0,
}

var (
	programs = make(map[innerapi.RenderProgramID]innerapi.RenderProgram)
)

func newTileRenderData(id innerapi.RenderDataID, pid innerapi.RenderProgramID, size api.BlockSize, slope api.BlockSlope, face api.Direction) *tileRenderData {
	log.Tracef("id: %s, pid: %s", id, pid)

	rrr := new(tileRenderData)
	rrr.id = id
	rrr.programId = pid
	
	switch slope {
		case api.SlopeNone:
			switch face {
				case api.DirectionSouth:
					rrr.instanceVertices = verticesSouth[:]
					rrr.instanceTextureUVs = uvsSouth[:]
				case api.DirectionNorth:
					rrr.instanceVertices = verticesNorth[:]
					rrr.instanceTextureUVs = uvsNorth[:]
				case api.DirectionEast:
					rrr.instanceVertices = verticesEast[:]
					rrr.instanceTextureUVs = uvsEast[:]
				case api.DirectionWest:
					rrr.instanceVertices = verticesWest[:]
					rrr.instanceTextureUVs = uvsWest[:]
				case api.DirectionUp:
					rrr.instanceVertices = verticesTop[:]
					rrr.instanceTextureUVs = uvsTop[:]
				case api.DirectionDown:
					rrr.instanceVertices = verticesBottom[:]
					rrr.instanceTextureUVs = uvsBottom[:]
				default:
					rrr.instanceVertices = []float32{}	
					rrr.instanceTextureUVs = []float32{}
			}
		case api.SlopeEastBottom:
			rrr.instanceVertices = verticesBottomEast[:]
			rrr.instanceTextureUVs = uvsEast[:]
		case api.SlopeEastTop:
			rrr.instanceVertices = verticesTopEast[:]
			rrr.instanceTextureUVs = uvsEast[:]
		case api.SlopeSouthBottom:
			rrr.instanceVertices = verticesBottomSouth[:]
			rrr.instanceTextureUVs = uvsSouth[:]
		case api.SlopeSouthEast:
			rrr.instanceVertices = verticesSouthEast[:]
			rrr.instanceTextureUVs = uvsSouth[:]
		case api.SlopeSouthTop:
			rrr.instanceVertices = verticesTopSouth[:]
			rrr.instanceTextureUVs = uvsSouth[:]
		case api.SlopeSouthWest:
			rrr.instanceVertices = verticesSouthWest[:]
			rrr.instanceTextureUVs = uvsSouth[:]
		case api.SlopeNorthBottom:
			rrr.instanceVertices = verticesBottomNorth[:]
			rrr.instanceTextureUVs = uvsNorth[:]
		case api.SlopeNorthEast:
			rrr.instanceVertices = verticesNorthEast[:]
			rrr.instanceTextureUVs = uvsNorth[:]
		case api.SlopeNorthTop:			
			rrr.instanceVertices = verticesTopNorth[:]
			rrr.instanceTextureUVs = uvsNorth[:]
		case api.SlopeNorthWest:
			rrr.instanceVertices = verticesNorthWest[:]
			rrr.instanceTextureUVs = uvsWest[:]
		case api.SlopeWestBottom:
			rrr.instanceVertices = verticesBottomWest[:]
			rrr.instanceTextureUVs = uvsWest[:]
		case api.SlopeWestTop:
			rrr.instanceVertices = verticesTopWest[:]
			rrr.instanceTextureUVs = uvsWest[:]
		default:
			rrr.instanceVertices = []float32{}
			rrr.instanceTextureUVs = []float32{}										
	}
	
	rrr.primitiveCount = len(rrr.instanceVertices) / 3
	
	tmp := make([]float32, len(rrr.instanceVertices))
	copy(tmp, rrr.instanceVertices)
	
	multiplier := size.Multiplier() 
	
	for i,c := range tmp {
		tmp[i] = c * multiplier
	}
	
	rrr.instanceVertices = tmp
	
	return rrr
}

func GetProgramForBlockType(pid innerapi.RenderProgramID, br modules.BlockRegistry) innerapi.RenderProgram {
	log.Tracef("%v", pid)
	p, ok := programs[pid]

	if !ok {
		p = innerapi.RenderProgram{}
		p.Program = gl.CreateProgram()
		var err error
		err = glutils.LoadShader(p.Program, gl.VERTEX_SHADER, "shaders/block.vertexshader")
		if err != nil {
			log.Errorf("Cannot load shader %#v: %#v", "shaders/block.vertexshader", err)
		}
		err = glutils.LoadShader(p.Program, gl.FRAGMENT_SHADER, "shaders/block.fragmentshader")
		if err != nil {
			log.Errorf("Cannot load shader %#v: %#v", "shaders/block.fragmentshader", err)
		}

		gl.LinkProgram(p.Program)

		var status int32
		gl.GetProgramiv(p.Program, gl.LINK_STATUS, &status)
		if status != gl.TRUE {
			var logLength int32
			gl.GetProgramiv(p.Program, gl.INFO_LOG_LENGTH, &logLength)

			lg := strings.Repeat("\x00", int(logLength+1))
			gl.GetProgramInfoLog(p.Program, logLength, nil, gl.Str(lg))

			log.Errorf("linker error: %v", lg)
		}

		p.MVP = gl.GetUniformLocation(p.Program, gl.Str("MVP\x00"))
		p.M = gl.GetUniformLocation(p.Program, gl.Str("M\x00"))
		p.V = gl.GetUniformLocation(p.Program, gl.Str("V\x00"))
		p.Texs = gl.GetUniformLocation(p.Program, gl.Str("textureSampler\x00"))
		p.Norm = gl.GetUniformLocation(p.Program, gl.Str("normalTextureSampler\x00"))
		p.LightDir = gl.GetUniformLocation(p.Program, gl.Str("lightDirection\x00"))
		p.SkyColor = gl.GetUniformLocation(p.Program, gl.Str("skyColor\x00"))
		p.GroundColor = gl.GetUniformLocation(p.Program, gl.Str("groundColor\x00"))
		p.Shift = gl.GetUniformLocation(p.Program, gl.Str("shift\x00"))
		
		p.Vertex = uint32(gl.GetAttribLocation(p.Program, gl.Str("vertexPosition_modelspace\x00")))
		p.Color = uint32(gl.GetAttribLocation(p.Program, gl.Str("vertexTextureUV\x00")))
		p.Normal = uint32(gl.GetAttribLocation(p.Program, gl.Str("vertexNormal_modelspace\x00")))
		p.Position = uint32(gl.GetAttribLocation(p.Program, gl.Str("primitivePosition_worldspace\x00")))
		
		id, _ := strconv.Atoi(string(pid))
		block := br.ByID(uint32(id))
		
		//log.Warnf("%v %v %s -> %s", pid, id, block.Type, br.Dir(block))
		
		texturePath := fmt.Sprintf("%s/texture.dds", br.Dir(block))
		p.Texture, err = glutils.CreateTextureFromDDS(texturePath)

		if err != nil {
			log.Errorf("Cannot load texture %#v: %#v", texturePath, err)
		}
		
		texturePath = fmt.Sprintf("%s/normal.dds", br.Dir(block))
		p.NormalTexture, err = glutils.CreateTextureFromDDS(texturePath)

		if err != nil {
			log.Errorf("Cannot load texture %#v: %#v", texturePath, err)
		}

		programs[pid] = p
	}

	log.Debugf("Program %#v", p)

	return p
}

type tileRenderData struct {
	id        innerapi.RenderDataID
	programId innerapi.RenderProgramID

	primitiveCount int
	instanceCount int
	
	vao uint32
	vbo uint32
	cbo uint32
	pbo uint32

	instanceVertices    []float32
	instanceTextureUVs []float32
	positions   []float32

	p innerapi.RenderProgram
}

func (this *tileRenderData) SetBlockPositions(positions []float32) {
	this.positions = positions
	this.instanceCount = len(positions) / 3
}

func (this *tileRenderData) RenderDataID() innerapi.RenderDataID {
	return this.id
}

func (this *tileRenderData) RenderProgramIDs() []innerapi.RenderProgramID {
	return []innerapi.RenderProgramID{this.programId}
}

func (this *tileRenderData) ID() string {
	return string(this.id)
}

func (this *tileRenderData) Target() string {
	return "graphics/gl"
}

func (this *tileRenderData) Init(p []innerapi.RenderProgram) {
	if this.instanceVertices != nil {
		gl.DeleteVertexArrays(1, &this.vao)
		gl.GenVertexArrays(1, &this.vao)
		
		gl.DeleteBuffers(1, &this.vbo)
		gl.GenBuffers(1, &this.vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(this.instanceVertices)*4, gl.Ptr(this.instanceVertices), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		this.instanceVertices = nil

		gl.DeleteBuffers(1, &this.cbo)
		gl.GenBuffers(1, &this.cbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, this.cbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(this.instanceTextureUVs)*4, gl.Ptr(this.instanceTextureUVs), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		this.instanceTextureUVs = nil
		
		if len(p) > 0 {
			this.p = p[0]
		} else {
			log.Errorf("No RenderProgramID provided for %v", this.programId)
		}
	}
	
	if this.positions != nil {
		gl.DeleteBuffers(1, &this.pbo)
		gl.GenBuffers(1, &this.pbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, this.pbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(this.positions)*4, gl.Ptr(this.positions), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		this.positions = nil
	}
}

func (this *tileRenderData) Render(camera innerapi.GlCamera) {
	gl.Enable(gl.CULL_FACE)
	gl.FrontFace(gl.CW)
	
	gl.BindVertexArray(this.vao)
	
	gl.UseProgram(this.p.Program)
	
	gl.UniformMatrix4fv(this.p.MVP, 1, false, &camera.MVP[0])
	gl.UniformMatrix4fv(this.p.M, 1, false, &camera.M[0])
	gl.UniformMatrix4fv(this.p.V, 1, false, &camera.V[0])
	
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, this.p.Texture)
	gl.Uniform1i(this.p.Texs, 0)
	
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, this.p.NormalTexture)
	gl.Uniform1i(this.p.Norm, 1)

	//gl.Uniform4fv(this.p.FogColor, 1, &FOGCOLOR[0])
	
	gl.EnableVertexAttribArray(this.p.Vertex)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.VertexAttribPointer(this.p.Vertex, 3, gl.FLOAT, false, 0, gl.Ptr(nil))
	
	gl.EnableVertexAttribArray(this.p.Color)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.cbo)
	gl.VertexAttribPointer(this.p.Color, 2, gl.FLOAT, false, 0, gl.Ptr(nil))
	
	gl.EnableVertexAttribArray(this.p.Position)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.pbo)
	gl.VertexAttribPointer(this.p.Position, 3, gl.FLOAT, false, 0, gl.Ptr(nil))
	gl.VertexAttribDivisor(this.p.Position, 1)
	
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, int32(this.primitiveCount), int32(this.instanceCount))
	
	gl.DisableVertexAttribArray(this.p.Vertex)
	gl.DisableVertexAttribArray(this.p.Color)
	gl.DisableVertexAttribArray(this.p.Position)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	
	gl.UseProgram(0)
	
	gl.BindVertexArray(0)
	
	gl.Disable(gl.CULL_FACE)
}

func (this *tileRenderData) Purge() {
	gl.DeleteVertexArrays(1, &this.vao)
	gl.DeleteBuffers(1, &this.vbo)
	gl.DeleteBuffers(1, &this.cbo)
	gl.DeleteBuffers(1, &this.pbo)
}
