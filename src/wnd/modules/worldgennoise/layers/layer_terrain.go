package layers

import (
	"wnd/api"
	"wnd/utils/log"

	"fmt"
	//"math"
	"sync"
)

const (
	_ELEVATION_CACHE_FORMAT = "%d/%d"
)

type TerrainLayer struct {
	anlGenerator
	sync.Mutex
	universe       *api.Universe
	blockCache     map[string]uint32
	elevationCache map[string]float64
}

func NewTerrainLayer(u *api.Universe, supportedBlocks map[string]api.BlockDefinition) *TerrainLayer {
	m := make(map[string]uint32)

	for k, v := range supportedBlocks {
		m[k] = v.ID
	}

	t := &TerrainLayer{
		/*perlinGenerator: perlinGenerator{
			octaves: func(seed int) int {
				return 4
			},
			scale: func(seed float64) float64 {
				return seed
			},
			seeder: func(seed int64) int64 {
				return seed / 71
			},
			alpha: func(seed float64) float64 {
				return seed
			},
			beta: func(seed float64) float64 {
				return seed
			},
		},*/
		universe:       u,
		blockCache:     m,
		elevationCache: make(map[string]float64),
	}

	t.init(u.Seed)
	return t
}

func (this *TerrainLayer) ID() string {
	return "TerrainLayer"
}

func (this *TerrainLayer) obtain(x, z int32, seed int64, chunk *api.Chunk) float64 {
	var vf float64

	xz := fmt.Sprintf(_ELEVATION_CACHE_FORMAT, x, z)
	if cached, ok := this.elevationCache[xz]; ok {
		//log.Debugf("TerrainLayer.Fill %d/%d/%d: got cached %v", x, z, chunk.Coords.Y, cached)
		vf = cached
	} else {
		vf = this.fill(
			int32(x),
			int32(z),
			seed,
			float64(seed % 14), //float64(chunk.BiomeData.HeightBias)-float64(chunk.BiomeData.Temperature),
			float64(seed % 16),//float64(chunk.BiomeData.Humidity)*float64(chunk.BiomeData.SeaLevelHeight),
			10, //3000/float64(chunk.BiomeData.HeightBias),
			6, //int(chunk.BiomeData.HeightBias)*int(chunk.BiomeData.SeaLevelHeight),
		)

		//vf /= math.MaxUint8
		//vf *= math.MaxInt8
		//vf -= math.MaxInt32 / 2

		this.Lock()
		this.elevationCache[xz] = vf
		this.Unlock()
		//e := float64(chunk.BiomeData.SeaLevelHeight) / float64(math.MaxInt8) * float64(math.MaxInt32 * 0.66)

		//log.Debugf("TerrainLayer.Fill %d/%d/%d: got %d / e %d sealevel %d", x, z, chunk.Coords.Y, intvf, int32(e), chunk.BiomeData.SeaLevelHeight)
	}

	return vf
}

func (this *TerrainLayer) Fill(seed int64, chunk *api.Chunk, size uint32) (e error) {
	log.Tracef("TerrainLayer.fill: seed %v, chunk %v (biome %#v)", seed, chunk.Coords.Pack(), chunk.BiomeData)

	chunk.PathDown = 0
	chunk.PathEast = 0
	chunk.PathSouth = 0
	chunk.PathNorth = 0
	chunk.PathUp = 0
	chunk.PathWest = 0

	for x := chunk.Coords.X; x < (chunk.Coords.X + int32(api.ChunkSideSize)); x++ {
		for z := chunk.Coords.Z; z < (chunk.Coords.Z + int32(api.ChunkSideSize)); z++ {

			vf := this.obtain(x, z, seed, chunk)
			intvf := int32(vf)

			vfe := int32(this.obtain(x-1, z, seed, chunk))
			vfw := int32(this.obtain(x+1, z, seed, chunk))
			vfn := int32(this.obtain(x, z+1, seed, chunk))
			vfs := int32(this.obtain(x, z-1, seed, chunk))

			vfne := int32(this.obtain(x-1, z+1, seed, chunk))
			vfnw := int32(this.obtain(x+1, z+1, seed, chunk))
			vfse := int32(this.obtain(x-1, z-1, seed, chunk))
			vfsw := int32(this.obtain(x+1, z-1, seed, chunk))

			//vf := int32((vf1 * math.MaxFloat64) / (math.MaxFloat64 / math.MaxInt32))

			//if vf-e < math.MaxInt16 {
			//vf = math.Min(e, vf) + (vf-e)
			//}

			//log.Warnf("TerrainLayer.Fill %d/%d/%d: got %d (%v)", x, z, chunk.Coords.Y, intvf, vf)

			for y := chunk.Coords.Y; y < (chunk.Coords.Y + int32(api.ChunkSideSize)); y++ {
				if intvf >= y {
					coords := api.WorldCoords{X: int32(x), Y: int32(y), Z: int32(z)}

					slope := api.SlopeNone

					if intvf == y {
						if intvf <= vfn && intvf > vfs {
							if intvf < vfn || (vfn <= vfne && vfn <= vfnw) {
								slope = api.SlopeSouthTop
							}
						} else if intvf <= vfs && intvf > vfn {
							if intvf < vfs || (vfs <= vfse && vfs <= vfsw) {
								slope = api.SlopeNorthTop
							}
						}

						if intvf <= vfe && intvf > vfw {
							if intvf < vfe || (vfe <= vfne && vfe <= vfse) {
								slope = api.SlopeWestTop
							}
						} else if intvf <= vfw && intvf > vfe {
							if intvf < vfw || (vfw <= vfnw && vfw <= vfsw) {
								slope = api.SlopeEastTop
							}
						}
					}

					if chunk.BlockAt(coords) == nil {
						if intvf == y && x == chunk.Coords.X && z == chunk.Coords.Z {
							mark := &api.Block{
								ID:          this.blockCache["tile"],
								Orientation: api.OrientationDefault,
								Size:        api.SizeQuarter,
								Slope:       api.SlopeNone,
							}
							
							markf := &api.Block{
								ID:          this.blockCache["grass"],
								Orientation: api.OrientationDefault,
								Size:        api.SizeQuarter,
								Slope:       api.SlopeNone,
							}

							mcoords := api.WorldCoords{X: int32(x), Y: int32(y + 1), Z: int32(z)}
							
							icoords1 := api.InnerCoords{X: 2, Y: 1, Z: 1}
							icoords2 := api.InnerCoords{X: 1, Y: 2, Z: 1}
							icoords3 := api.InnerCoords{X: 3, Y: 2, Z: 1}
							icoords4 := api.InnerCoords{X: 2, Y: 2, Z: 0}
							icoords5 := api.InnerCoords{X: 2, Y: 2, Z: 2}
							icoords6 := api.InnerCoords{X: 2, Y: 3, Z: 1}
							icoords7 := api.InnerCoords{X: 2, Y: 0, Z: 1}

							chunk.SetBlockInner(mcoords, icoords1, markf)
							chunk.SetBlockInner(mcoords, icoords2, mark)
							chunk.SetBlockInner(mcoords, icoords3, mark)
							chunk.SetBlockInner(mcoords, icoords4, mark)
							chunk.SetBlockInner(mcoords, icoords5, mark)
							chunk.SetBlockInner(mcoords, icoords6, mark)
							chunk.SetBlockInner(mcoords, icoords7, markf)
						}

						block := &api.Block{
							ID:          this.blockCache["cobblestone"],
							Orientation: api.OrientationDefault,
							Size:        api.SizeFull,
							Slope:       slope,
						}

						chunk.SetBlock(coords, block)
					}

					blockInChunk := chunk.BlockAt(coords)

					blockInChunk.Adjacents = blockInChunk.Adjacents.Add(api.DirectionDown)

					if intvf > y {
						blockInChunk.Adjacents = blockInChunk.Adjacents.Add(api.DirectionUp)
					} else {
						chunk.PathUp++
					}

					if intvf < chunk.Coords.Y {
						chunk.PathDown++
					}

					//if x == chunk.Coords.X {
					easty := int32(vfe)
					if x-1 >= -int32(size) && easty > y {
						blockInChunk.Adjacents = blockInChunk.Adjacents.Add(api.DirectionEast)
					} else {
						chunk.PathEast++
					}
					//} else if x == (chunk.Coords.X + api.ChunkSideSize - 1) {
					westy := int32(vfw)
					if x+1 < int32(size) && westy > y {
						blockInChunk.Adjacents = blockInChunk.Adjacents.Add(api.DirectionWest)
					} else {
						chunk.PathWest++
					}
					//}

					//if z == chunk.Coords.Z {
					northy := int32(vfn)
					if z+1 < int32(size) && northy > y {
						blockInChunk.Adjacents = blockInChunk.Adjacents.Add(api.DirectionNorth)
					} else {
						chunk.PathSouth++
					}
					//} else if z == (chunk.Coords.Z + api.ChunkSideSize - 1) {
					southy := int32(vfs)
					if z-1 >= -int32(size) && southy > y {
						blockInChunk.Adjacents = blockInChunk.Adjacents.Add(api.DirectionSouth)
					} else {
						chunk.PathNorth++
					}
					//}
				}
			}
		}
	}
	
	//log.Warnf("tiles %v => u %v / d %v / e %v / w %v / n %v / s %v", chunk.Coords, chunk.PathUp, chunk.PathDown, chunk.PathEast, chunk.PathWest, chunk.PathSouth, chunk.PathNorth)

	return nil
}
