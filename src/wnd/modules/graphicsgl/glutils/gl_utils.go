package glutils

import (
	"wnd/api"
	"wnd/resources"
	"wnd/utils/log"
	"wnd/modules/graphicsgl/innerapi"

	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"sync"
	"io"
	"math"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/spate/glimage"
	
	glm "github.com/go-gl/mathgl/mgl32"
	ddstypes "github.com/spate/glimage/dds/types"
)

var (
	_DIAMETER = math.Sqrt(3 * api.ChunkSideSizeFloat64 * api.ChunkSideSizeFloat64 * 2)
)

type OrderedMap struct {
	sync.RWMutex
	
	n int
	k []innerapi.RenderDataID
	m map[innerapi.RenderDataID]innerapi.RenderData
}

func NewOrderedMap() (m *OrderedMap) {
	m = new(OrderedMap)
	m.Clear()
	return
}

func (this *OrderedMap) Put(d innerapi.RenderData) {
	if d == nil {
		return
	}
	
	this.Lock()
	defer this.Unlock()
	
	this.m[d.RenderDataID()] = d
	this.k = append(this.k, d.RenderDataID())
	this.n++
}

func (this *OrderedMap) Get(k innerapi.RenderDataID) (innerapi.RenderData, bool) {
	i,ok := this.m[k]
	return i,ok
}

func (this *OrderedMap) Delete(k innerapi.RenderDataID) (r innerapi.RenderData) {
	this.Lock()
	defer this.Unlock()
	
	ok := false
	
	if r,ok = this.m[k]; ok {
		delete(this.m, k)
		
		index := -1 
		
		for i,ik := range this.k {
			if k == ik {
				index = i
				break
			}
		}
		
		if index > -1 {
			this.k = append(this.k[:index], this.k[index+1:]...)
		}
		
		this.n--
	} else {
		r = nil
	}
	
	return
}

func (this *OrderedMap) Keys() []innerapi.RenderDataID {
	return this.k[:this.n]
}

func (this *OrderedMap) Size() int {
	return this.n
}

func (this *OrderedMap) Clear() {
	this.n = 0
	this.k = make([]innerapi.RenderDataID, 0, api.ChunkSideSize * api.ChunkSideSize * api.ChunkSideSize)
	this.m = make(map[innerapi.RenderDataID]innerapi.RenderData)
}

func glError() {
	log.Tracef("GL error: %#v", gl.GetError())
}

func LoadShader(program, shaderType uint32, shaderName string) error {
	shader := gl.CreateShader(shaderType)
	defer gl.DeleteShader(shader)

	shaderSrc, e := resources.Asset(shaderName)
	if e != nil {
		return e
	}

	cs,free := gl.Strs(string(shaderSrc) + "\x00")
	gl.ShaderSource(shader, 1, cs, nil)
	free();
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return errors.New(log)
	}

	gl.AttachShader(program, shader)
	return nil
}

func CreateTextureFromDDS(filename string) (uint32, error) {
	src, err := os.Open(filename)
	if err != nil {
		log.Warnf("Cannot load texture: %#v", filename)
		return 0, err
	}

	var r io.Reader
	var header ddstypes.DDS_HEADER
	var tmp [128]byte

	r = bufio.NewReader(src)
	// Check for DDS magic number
	_, err = io.ReadFull(r, tmp[:4])
	if err != nil {
		return 0, err
	}
	ident := string(tmp[0:4])
	if ident != "DDS " {
		return 0, fmt.Errorf("dds: wrong magic number")
	}
	// Decode the DDS header
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return 0, err
	}
	if header.Size != 124 {
		return 0, fmt.Errorf("dds: invalid DDS header")
	}
	if header.Ddspf.FourCC == ddstypes.FOURCC_DX10 {
		return 0, fmt.Errorf("dds: unsupported DX10 header")
	}

	// Check if it's a supported format
	// For now, we'll only support DXT1,DXT3,DXT5
	neededFlags := uint32(ddstypes.DDSD_HEIGHT | ddstypes.DDSD_WIDTH | ddstypes.DDSD_PIXELFORMAT)
	if header.Flags&neededFlags != neededFlags {
		return 0, fmt.Errorf("dds: file header is missing necessary dds flags")
	}
	// Sanitize mipmap count
	if header.Flags&ddstypes.DDSD_MIPMAPCOUNT < 1 {
		header.MipMapCount = 1
	}

	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	switch {
	case header.Ddspf.Flags&ddstypes.DDPF_FOURCC != 0:
		switch header.Ddspf.FourCC {
		case ddstypes.FOURCC_DXT1:
			w, h := int(header.Width), int(header.Height)
			for i := 0; i < int(header.MipMapCount); i++ {
				log.Tracef("mipmap %v is %vx%v", i, w, h)
				img := glimage.NewDxt1(image.Rect(0, 0, w, h))
				_, err = io.ReadFull(r, img.Pix)
				if err != nil {
					return 0, err
				}
				gl.CompressedTexImage2D(gl.TEXTURE_2D, int32(i), 0x83f1 /*gl.COMPRESSED_RGBA_S3TC_DXT1_EXT*/, int32(w), int32(h), 0, int32(len(img.Pix)), gl.Ptr(img.Pix))
				w >>= 1
				h >>= 1

				if w < 1 {
					w = 1
				}
				if h < 1 {
					h = 1
				}
			}
		case ddstypes.FOURCC_DXT3:
			w, h := int(header.Width), int(header.Height)
			for i := 0; i < int(header.MipMapCount); i++ {
				log.Tracef("mipmap %v is %vx%v", i, w, h)
				img := glimage.NewDxt3(image.Rect(0, 0, w, h))
				_, err = io.ReadFull(r, img.Pix)
				if err != nil {
					return 0, err
				}
				gl.CompressedTexImage2D(gl.TEXTURE_2D, int32(i), 0x83f2 /*gl.COMPRESSED_RGBA_S3TC_DXT3_EXT*/, int32(w), int32(h), 0, int32(len(img.Pix)), gl.Ptr(img.Pix))
				w >>= 1
				h >>= 1
				if w < 1 {
					w = 1
				}
				if h < 1 {
					h = 1
				}
			}
		case ddstypes.FOURCC_DXT5:
			w, h := int(header.Width), int(header.Height)
			for i := 0; i < int(header.MipMapCount); i++ {
				log.Tracef("mipmap %v is %vx%v", i, w, h)
				img := glimage.NewDxt5(image.Rect(0, 0, w, h))
				_, err = io.ReadFull(r, img.Pix)
				if err != nil {
					return 0, err
				}
				gl.CompressedTexImage2D(gl.TEXTURE_2D, int32(i), 0x83f3 /*gl.COMPRESSED_RGBA_S3TC_DXT5_EXT*/, int32(w), int32(h), 0, int32(len(img.Pix)), gl.Ptr(img.Pix))
				w >>= 1
				h >>= 1
				if w < 1 {
					w = 1
				}
				if h < 1 {
					h = 1
				}
			}
		default:
			return 0, fmt.Errorf("dds: unrecognized format %v", header.Ddspf)
		}
	case header.Ddspf.Flags&ddstypes.DDPF_RGB != 0:
		// Color formats
		if header.Ddspf.Flags&ddstypes.DDPF_ALPHAPIXELS != 0 {
			// Color formats with alpha
			switch {
			// A8R8G8B8
			case header.Ddspf.RBitMask == 0x00FF0000 && header.Ddspf.GBitMask == 0x0000FF00 &&
				header.Ddspf.BBitMask == 0x000000FF && header.Ddspf.ABitMask == 0xFF000000:
				w, h := int(header.Width), int(header.Height)
				for i := 0; i < int(header.MipMapCount); i++ {
					log.Tracef("mipmap %v is %vx%v", i, w, h)
					img := glimage.NewBGRA(image.Rect(0, 0, w, h))
					_, err = io.ReadFull(r, img.Pix)
					if err != nil {
						return 0, err
					}
					gl.TexImage2D(gl.TEXTURE_2D, int32(i), gl.RGBA, int32(w), int32(h), 0, gl.BGRA, gl.UNSIGNED_INT_8_8_8_8, gl.Ptr(img.Pix))
					w >>= 1
					h >>= 1
					if w < 1 {
						w = 1
					}
					if h < 1 {
						h = 1
					}
				}
			// A4R4G4B4
			case header.Ddspf.RBitMask == 0x0F00 && header.Ddspf.GBitMask == 0x00F0 &&
				header.Ddspf.BBitMask == 0x000F && header.Ddspf.ABitMask == 0xF000:
				w, h := int(header.Width), int(header.Height)
				for i := 0; i < int(header.MipMapCount); i++ {
					log.Tracef("mipmap %v is %vx%v", i, w, h)
					img := glimage.NewBGRA4444(image.Rect(0, 0, w, h))
					err = binary.Read(r, binary.LittleEndian, &img.Pix)
					if err != nil {
						return 0, err
					}
					gl.TexImage2D(gl.TEXTURE_2D, int32(i), gl.RGBA, int32(w), int32(h), 0, gl.BGRA, gl.UNSIGNED_SHORT_4_4_4_4, gl.Ptr(img.Pix))
					w >>= 1
					h >>= 1
					if w < 1 {
						w = 1
					}
					if h < 1 {
						h = 1
					}
				}
			// A1R5G5B5
			case header.Ddspf.RBitMask == 0x7C00 && header.Ddspf.GBitMask == 0x03E0 &&
				header.Ddspf.BBitMask == 0x001F && header.Ddspf.ABitMask == 0x8000:
				w, h := int(header.Width), int(header.Height)
				for i := 0; i < int(header.MipMapCount); i++ {
					log.Tracef("mipmap %v is %vx%v", i, w, h)
					img := glimage.NewBGRA5551(image.Rect(0, 0, w, h))
					err = binary.Read(r, binary.LittleEndian, &img.Pix)
					if err != nil {
						return 0, err
					}
					gl.TexImage2D(gl.TEXTURE_2D, int32(i), gl.RGBA, int32(w), int32(h), 0, gl.BGRA, gl.UNSIGNED_SHORT_5_5_5_1, gl.Ptr(img.Pix))
					w >>= 1
					h >>= 1
					if w < 1 {
						w = 1
					}
					if h < 1 {
						h = 1
					}
				}
			default:
				return 0, fmt.Errorf("dds: unrecognized format %v", header.Ddspf)
			}
		} else {
			// Color formats without alpha
			switch {
			// R5G6B5
			case header.Ddspf.RBitMask == 0xF800 && header.Ddspf.GBitMask == 0x07E0 &&
				header.Ddspf.BBitMask == 0x001F && header.Ddspf.ABitMask == 0x0000:
				w, h := int(header.Width), int(header.Height)
				for i := 0; i < int(header.MipMapCount); i++ {
					log.Tracef("mipmap %v is %vx%v", i, w, h)
					img := glimage.NewBGR565(image.Rect(0, 0, w, h))
					err = binary.Read(r, binary.LittleEndian, &img.Pix)
					if err != nil {
						return 0, err
					}
					gl.TexImage2D(gl.TEXTURE_2D, int32(i), gl.RGB, int32(w), int32(h), 0, gl.BGR, gl.UNSIGNED_SHORT_5_6_5, gl.Ptr(img.Pix))
					w >>= 1
					h >>= 1
					if w < 1 {
						w = 1
					}
					if h < 1 {
						h = 1
					}
				}
			default:
				return 0, fmt.Errorf("dds: unrecognized format %v", header.Ddspf)
			}
		}
	default:
		return 0, fmt.Errorf("dds: unrecognized format %v", header.Ddspf)
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	return tex, nil
}

func IsChunkInFrustum(chunk *api.Chunk, mvp glm.Mat4) bool {
	crds := chunk.Coords.ToCoords()
	clip := glm.Mat4(mvp).Mul4x1(glm.Vec4{crds.X + float32(api.ChunkMiddleBlock), crds.Y + float32(api.ChunkMiddleBlock), crds.Z + float32(api.ChunkMiddleBlock), 1.0})
	clip[0] /= clip[3]
	clip[1] /= clip[3]

	diameter := _DIAMETER

	if clip[2] < -float32(diameter) {
		return false
	}

	diameter /= math.Abs(float64(clip[3]))

	if math.Abs(float64(clip[0])) > 1+diameter || math.Abs(float64(clip[1])) > 1+diameter {
		return false
	} else {
		return true
	}
}

func IsInFrustum(coords api.WorldCoords, mvp glm.Mat4) bool {
	log.Debugf("%v for %v ...", coords, mvp)

	crds := coords.ToCoords()
	clip := glm.Mat4(mvp).Mul4x1(glm.Vec4{crds.X, crds.Y, crds.Z, 1.0})

	result := float32(math.Abs(float64(clip[0]))) < clip[3] && float32(math.Abs(float64(clip[1]))) < clip[3] && float32(math.Abs(float64(clip[2]))) < clip[3]

	log.Debugf(" ... %#v", result)

	return result
}

func OppositeDir(dir api.Direction) api.Direction {
	switch dir {
	case api.DirectionNorth:
		return api.DirectionSouth
	case api.DirectionSouth:
		return api.DirectionNorth
	case api.DirectionEast:
		return api.DirectionWest
	case api.DirectionWest:
		return api.DirectionEast
	case api.DirectionUp:
		return api.DirectionDown
	default:
		return api.DirectionUp
	}
}
