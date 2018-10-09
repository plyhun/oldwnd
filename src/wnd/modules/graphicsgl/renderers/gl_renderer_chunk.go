package renderers

import (
	"wnd/api"
	"wnd/base/graphics"
	"wnd/modules/graphicsgl/glutils"
	"wnd/modules/graphicsgl/innerapi"
	"wnd/utils/log"

	"strconv"
)

const useTileInstancing = false

type innerCoords struct {
	x uint8
	y uint8
	z uint8

	smallBlocks []innerCoords
}

type chunkRenderData struct {
	*api.Chunk

	visibleBlocks []innerCoords
	blockDatas    []*blockRenderData
	
	shift [3]float32
}

func CreateChunkRenderData(c *api.Chunk, w *api.World) innerapi.RenderData {
	if c == nil {
		return nil
	}

	crd := new(chunkRenderData)
	crd.Chunk = c

	crd.fillVisibles(w.Chunks.Data)
	
	if len(crd.visibleBlocks) < 1 {
		return nil
	}

	crd.fillBlockRenderData(w)

	return crd
}

func (this *chunkRenderData) Init(programs map[innerapi.RenderProgramID]innerapi.RenderProgram) {
	for _, b := range this.blockDatas {
		b.Init([]innerapi.RenderProgram{programs[b.RenderProgramIDs()[0]]})
	}
}

func (this *chunkRenderData) RenderProgramIDs() []innerapi.RenderProgramID {
	ret := make([]innerapi.RenderProgramID, len(this.blockDatas))

	for i, b := range this.blockDatas {
		ret[i] = b.RenderProgramIDs()[0]
	}

	return ret
}

func (this *chunkRenderData) findBlockRenderData(key innerapi.RenderDataID) *blockRenderData {
	if this.blockDatas == nil {
		return nil
	}

	for i, b := range this.blockDatas {
		if b.RenderDataID() == key {
			return this.blockDatas[i]
		}
	}

	return nil
}

func (this *chunkRenderData) fillBlockRenderData(w *api.World) {
	//this.Purge()

	for _, cc := range this.visibleBlocks {
		b := this.Chunk.Blocks[cc.x][cc.y][cc.z]

		if b == nil {
			continue
		}

		if useTileInstancing {
			/*for _,face := range b.Adjacents.Inverse().Parse() {
				key := innerapi.RenderDataID(string(b.PackWithoutCoords()) + " " + strconv.Itoa(int(face)))

				brd := this.findBlockRenderData(key)

				if brd == nil {
					brd = newTileRenderData(key, innerapi.RenderProgramID(strconv.Itoa(int(b.ID))), b.Size, b.Slope, face)

					this.blockDatas = append(this.blockDatas, brd)
				}

				brd.positions = append(brd.positions, float32(this.Chunk.Coords.X+int32(cc.x)), float32(this.Chunk.Coords.Y+int32(cc.y)), float32(this.Chunk.Coords.Z+int32(cc.z)))
				brd.instanceCount++
			}*/
		} else {
			if cc.smallBlocks == nil {
				key := innerapi.RenderDataID(string(b.PackWithoutCoords()))

				brd := this.findBlockRenderData(key)

				if brd == nil {
					brd = newBlockRenderData(key, innerapi.RenderProgramID(strconv.Itoa(int(b.ID))), b.Size, b.Slope)

					this.blockDatas = append(this.blockDatas, brd)
				}

				brd.positions.data = append(brd.positions.data, float32(this.Chunk.Coords.X+int32(cc.x)), float32(this.Chunk.Coords.Y+int32(cc.y)), float32(this.Chunk.Coords.Z+int32(cc.z)))
				brd.instanceCount++
			} else if b.SmallBlocks != nil {
				for _, ccc := range cc.smallBlocks {
					bb := b.SmallBlocks[ccc.x][ccc.y][ccc.z]

					if bb == nil {
						continue
					}

					key := innerapi.RenderDataID(string(bb.PackWithoutCoords()))

					brd := this.findBlockRenderData(key)

					if brd == nil {
						brd = newBlockRenderData(key, innerapi.RenderProgramID(strconv.Itoa(int(bb.ID))), bb.Size, bb.Slope)

						this.blockDatas = append(this.blockDatas, brd)
					}

					brd.positions.data = append(brd.positions.data, float32(this.Chunk.Coords.X+int32(cc.x))+(float32(ccc.x)/4), float32(this.Chunk.Coords.Y+int32(cc.y))+(float32(ccc.y)/4), float32(this.Chunk.Coords.Z+int32(cc.z))+(float32(ccc.z)/4))
					brd.instanceCount++
				}
			}
		}
	}
}

func (this *chunkRenderData) fillVisibles(otherChunks map[api.PackedWorldCoords]*api.Chunk) {
	this.visibleBlocks = make([]innerCoords, 0)

	log.Tracef("%v", this.Coords)

	for x, bx := range this.Chunk.Blocks {
		for y, by := range bx {
			for z, b := range by {
				if b == nil || b.Adjacents.HasAllAdjacents() {
					continue
				}

				inn := innerCoords{x: uint8(x), y: uint8(y), z: uint8(z)}

				if b.Block != nil {
					this.visibleBlocks = append(this.visibleBlocks, inn)
				} else if b.SmallBlocks != nil {
					inn.smallBlocks = make([]innerCoords, 0, 4*4*4)

					for ix, ibx := range b.SmallBlocks {
						for iy, iby := range ibx {
							for iz, bb := range iby {
								if bb == nil || bb.Adjacents.HasAllAdjacents() {
									continue
								}

								inn.smallBlocks = append(inn.smallBlocks, innerCoords{x: uint8(ix), y: uint8(iy), z: uint8(iz)})
							}
						}
					}

					if len(inn.smallBlocks) > 0 {
						this.visibleBlocks = append(this.visibleBlocks, inn)
					}
				}
			}
		}
	}
}

func (this *chunkRenderData) ID() string {
	return "chunk " + string(this.Chunk.Coords.Pack())
}

func (this *chunkRenderData) Target() string {
	return graphics.TARGET
}

func (this *chunkRenderData) RenderDataID() innerapi.RenderDataID {
	return innerapi.RenderDataID(graphics.GetGraphicsOutputableID(this.Chunk))
}

func (this *chunkRenderData) Render(camera *innerapi.GlCamera) {
	for _, v := range this.blockDatas {
		if v != nil {
			if this.shift[0] != 0 || this.shift[1] != 0 || this.shift[2] != 0 {
				v.shift = this.shift
				v.Render(camera)
				v.shift = [3]float32{0,0,0}
			}
			v.Render(camera)
		}
	}
}

func (this *chunkRenderData) IsVisible(camera *innerapi.GlCamera, w *api.World) bool {
	/*if this.Chunk.PathDown == 0 && this.Chunk.PathEast == 0 && this.Chunk.PathNorth == 0 && this.Chunk.PathSouth == 0 && this.Chunk.PathUp == 0 && this.Chunk.PathWest == 0 {
		return false
	}  */

	observerChunk := camera.Position.ToWorldCoords()

	if observerChunk.Equals(this.Chunk.Coords) {
		return true
	}
	
	if shift,ok := w.GetShift(this.Chunk.Coords.ToCoords(), camera.Observer.Position); ok {
		this.shift = [3]float32{shift.X, shift.Y, shift.Z}
		return true
	}

	//chunkDelta := api.WorldCoords{X: this.Chunk.Coords.X - observerChunk.X, Y: this.Chunk.Coords.Y - observerChunk.Y, Z: this.Chunk.Coords.Z - observerChunk.Z}
	//viewDelta := api.Coords{X: camera.LooksAt.X - camera.Position.X, Y: camera.LooksAt.Y - camera.Position.Y, Z: camera.LooksAt.Z - camera.Position.Z}

	return glutils.IsChunkInFrustum(this.Chunk, camera.MVP)
}

func (this *chunkRenderData) Purge() {
	if this.blockDatas != nil {
		for _, v := range this.blockDatas {
			if v != nil {
				v.Purge()
			}
		}
	}

	this.blockDatas = make([]*blockRenderData, 0)
}
