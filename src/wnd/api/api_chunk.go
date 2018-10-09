package api

import (
	"wnd/utils/log"
	
	"sync"
)

/*type VisibilityStatus uint8

const (
	VisibilityStatusUnknown VisibilityStatus = iota
	VisibilityStatusVisible
	VisibilityStatusInvisible
)*/

type InnerCoords struct {
	X uint8
	Y uint8
	Z uint8
}

type BlockInChunk struct {
	*Block
	*SmallBlocks
	Tileable
}

type SmallBlocks [4][4][4]*BlockInChunk

type Chunk struct {
	BiomeData
	sync.RWMutex
	
	Coords WorldCoords
	Blocks [ChunkSideSize][ChunkSideSize][ChunkSideSize]*BlockInChunk
	
	PathNorth, PathSouth, PathEast, PathWest, PathUp, PathDown uint8
	
	VisibleBlocks []*InnerCoords
	
}

func (this *Chunk) CoordsToChunkCoords(coords WorldCoords) (x, y, z int32, ok bool) {
	x, y, z = coords.X-this.Coords.X, coords.Y-this.Coords.Y, coords.Z-this.Coords.Z

	switch {
	case x < 0, y < 0, z < 0, x >= int32(ChunkSideSize), y >= int32(ChunkSideSize), z >= int32(ChunkSideSize):
		ok = false
	default:
		ok = true
	}
	
	//log.Tracef("Chunk.CoordsToChunkCoords: %v - %v = %v, %v, %v (%v)", coords, this.Coords, x, y, z, ok)
	
	return
}

func (this *Chunk) BlockAt(coords WorldCoords) *BlockInChunk {
	//log.Tracef("Chunk.BlockAt: %#v", coords)

	//packed := coords.Pack()
	
	x,y,z,ok := this.CoordsToChunkCoords(coords)
	if !ok {
		return nil
	}

	return this.Blocks[x][y][z]
}

func (this *Chunk) SetBlock(coords WorldCoords, block *Block) {
	//log.Tracef("Chunk.SetBlock: %#v", block)

	x,y,z,ok := this.CoordsToChunkCoords(coords)
	if !ok {
		log.Infof("Cannot insert block %#v to chunk %#v", coords, this.Coords)
		return 
	}
	
	this.Lock()
	
	bic := &BlockInChunk{
		Block: block,
	}
	
	bic.Adjacents = bic.Adjacents.SetDefined(true)

	this.Blocks[x][y][z] = bic
	
	if x > 0 && this.Blocks[x-1][y][z] != nil && this.Blocks[x-1][y][z].Block != nil && this.Blocks[x-1][y][z].Slope == SlopeNone {/*(
			this.Blocks[x-1][y][z].Slope == SlopeNone || 
			this.Blocks[x-1][y][z].Slope == SlopeNorthWest ||
			this.Blocks[x-1][y][z].Slope == SlopeSouthWest || 
			this.Blocks[x-1][y][z].Slope == SlopeWestTop ||
			this.Blocks[x-1][y][z].Slope == SlopeWestBottom) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionEast)
		
		if bic.Slope == SlopeNone {/*|| 
			bic.Slope == SlopeNorthEast ||
			bic.Slope == SlopeSouthEast || 
			bic.Slope == SlopeEastTop ||
			bic.Slope == SlopeEastBottom {*/
			this.Blocks[x-1][y][z].Adjacents = this.Blocks[x-1][y][z].Adjacents.Add(DirectionWest)
		} else {
			this.Blocks[x-1][y][z].Adjacents = this.Blocks[x-1][y][z].Adjacents.Remove(DirectionWest)
		}
	} 
	
	if x < int32(ChunkSideSize)-1 && this.Blocks[x+1][y][z] != nil && this.Blocks[x+1][y][z].Block != nil && this.Blocks[x+1][y][z].Slope == SlopeNone {/*(
			this.Blocks[x+1][y][z].Slope == SlopeNone || 
			this.Blocks[x+1][y][z].Slope == SlopeNorthEast ||
			this.Blocks[x+1][y][z].Slope == SlopeSouthEast || 
			this.Blocks[x+1][y][z].Slope == SlopeEastTop ||
			this.Blocks[x+1][y][z].Slope == SlopeEastBottom) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionWest)
		
		if bic.Slope == SlopeNone {/*|| 
			bic.Slope == SlopeNorthWest ||
			bic.Slope == SlopeSouthWest || 
			bic.Slope == SlopeWestTop ||
			bic.Slope == SlopeWestBottom {*/
			this.Blocks[x+1][y][z].Adjacents = this.Blocks[x+1][y][z].Adjacents.Add(DirectionEast)
		} else {
			this.Blocks[x+1][y][z].Adjacents = this.Blocks[x+1][y][z].Adjacents.Remove(DirectionEast)
		}
	} 
	
	if y > 0 && this.Blocks[x][y-1][z] != nil && this.Blocks[x][y-1][z].Block != nil && this.Blocks[x][y-1][z].Slope == SlopeNone {/*(
			this.Blocks[x][y-1][z].Slope == SlopeNone || 
			this.Blocks[x][y-1][z].Slope == SlopeNorthTop ||
			this.Blocks[x][y-1][z].Slope == SlopeSouthTop || 
			this.Blocks[x][y-1][z].Slope == SlopeEastTop ||
			this.Blocks[x][y-1][z].Slope == SlopeWestTop) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionDown)
		
		if bic.Slope == SlopeNone {/*|| 
			bic.Slope == SlopeNorthBottom ||
			bic.Slope == SlopeSouthBottom || 
			bic.Slope == SlopeEastBottom ||
			bic.Slope == SlopeWestBottom {*/
			this.Blocks[x][y-1][z].Adjacents = this.Blocks[x][y-1][z].Adjacents.Add(DirectionUp)
		} else {
			this.Blocks[x][y-1][z].Adjacents = this.Blocks[x][y-1][z].Adjacents.Remove(DirectionUp)
		}
	} 
	
	if y < int32(ChunkSideSize)-1 && this.Blocks[x][y+1][z] != nil && this.Blocks[x][y+1][z].Block != nil && this.Blocks[x][y+1][z].Block != nil && this.Blocks[x][y+1][z].Slope == SlopeNone {/*(
			this.Blocks[x][y+1][z].Slope == SlopeNone || 
			this.Blocks[x][y+1][z].Slope == SlopeNorthBottom ||
			this.Blocks[x][y+1][z].Slope == SlopeSouthBottom || 
			this.Blocks[x][y+1][z].Slope == SlopeEastBottom ||
			this.Blocks[x][y+1][z].Slope == SlopeWestBottom) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionUp)
		
		if bic.Slope == SlopeNone {/*|| 
			bic.Slope == SlopeNorthTop ||
			bic.Slope == SlopeSouthTop || 
			bic.Slope == SlopeEastTop ||
			bic.Slope == SlopeWestTop {*/
			this.Blocks[x][y+1][z].Adjacents = this.Blocks[x][y+1][z].Adjacents.Add(DirectionDown)
		} else {
			this.Blocks[x][y+1][z].Adjacents = this.Blocks[x][y+1][z].Adjacents.Remove(DirectionDown)
		}
	} 
	
	if z > 0 && this.Blocks[x][y][z-1] != nil && this.Blocks[x][y][z-1].Block != nil && this.Blocks[x][y][z-1].Slope == SlopeNone {/*(
			this.Blocks[x][y][z-1].Slope == SlopeNone || 
			this.Blocks[x][y][z-1].Slope == SlopeSouthBottom ||
			this.Blocks[x][y][z-1].Slope == SlopeSouthTop || 
			this.Blocks[x][y][z-1].Slope == SlopeSouthEast ||
			this.Blocks[x][y][z-1].Slope == SlopeSouthWest) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionSouth)
		
		if bic.Slope == SlopeNone {/* || 
			bic.Slope == SlopeNorthTop ||
			bic.Slope == SlopeNorthBottom || 
			bic.Slope == SlopeNorthEast ||
			bic.Slope == SlopeNorthWest {*/
			this.Blocks[x][y][z-1].Adjacents = this.Blocks[x][y][z-1].Adjacents.Add(DirectionNorth)
		} else {
			this.Blocks[x][y][z-1].Adjacents = this.Blocks[x][y][z-1].Adjacents.Remove(DirectionNorth)
		}
	} 
	
	if z < int32(ChunkSideSize)-1 && this.Blocks[x][y][z+1] != nil && this.Blocks[x][y][z+1].Block != nil && this.Blocks[x][y][z+1].Slope == SlopeNone {/*(
			this.Blocks[x][y][z+1].Slope == SlopeNone || 
			this.Blocks[x][y][z+1].Slope == SlopeNorthBottom ||
			this.Blocks[x][y][z+1].Slope == SlopeNorthTop || 
			this.Blocks[x][y][z+1].Slope == SlopeNorthEast ||
			this.Blocks[x][y][z+1].Slope == SlopeNorthWest) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionNorth)
		
		if bic.Slope == SlopeNone {/*|| 
			bic.Slope == SlopeSouthTop ||
			bic.Slope == SlopeSouthBottom || 
			bic.Slope == SlopeSouthEast ||
			bic.Slope == SlopeSouthWest {*/
			this.Blocks[x][y][z+1].Adjacents = this.Blocks[x][y][z+1].Adjacents.Add(DirectionSouth)
		} else {
			this.Blocks[x][y][z+1].Adjacents = this.Blocks[x][y][z+1].Adjacents.Remove(DirectionSouth)
		}
	} 
			
	//this.addDirty(block.PackWithoutCoords())
	
	this.Unlock()
}

func (this *Chunk) BlockAtInner(coords WorldCoords, inner InnerCoords) *BlockInChunk {
	//log.Tracef("Chunk.BlockAt: %#v", coords)

	container := this.BlockAt(coords)
	
	if container == nil || container.SmallBlocks == nil || inner.X > 3 || inner.Y > 3 || inner.Z > 3 {
		return nil
	}

	return container.SmallBlocks[inner.X][inner.Y][inner.Z]
}

func (this *Chunk) SetBlockInner(coords WorldCoords, inner InnerCoords, block *Block) {
	//log.Tracef("Chunk.SetBlock: %#v", block)

	x,y,z,ok := this.CoordsToChunkCoords(coords)
	if !ok {
		log.Infof("Cannot insert block %#v to chunk %#v", coords, this.Coords)
		return 
	}
	
	this.Lock()
	
	container := this.BlockAt(coords)
	
	if (container != nil && container.Block != nil) || inner.X > 3 || inner.Y > 3 || inner.Z > 3 || (container != nil && container.SmallBlocks != nil && container.SmallBlocks[inner.X][inner.Y][inner.Z] != nil){
		return 
	}
	
	if block.Size < SizeHalf {
		this.SetBlock(coords, block)
	}
	
	if container == nil {
		container = &BlockInChunk{}
		this.Blocks[x][y][z] = container
	}
	
	if container.SmallBlocks == nil {
		var sb SmallBlocks
		container.SmallBlocks = &sb
	}
	
	bic := &BlockInChunk{
		Block: block,
	}
	
	bic.Adjacents = bic.Adjacents.SetDefined(true)

	if block.Size == SizeHalf {
		if inner.X > 2 || inner.Y > 2 || inner.Z > 2 ||
				container.SmallBlocks[inner.X+1][inner.Y][inner.Z] != nil ||
				container.SmallBlocks[inner.X][inner.Y+1][inner.Z] != nil ||
				container.SmallBlocks[inner.X][inner.Y][inner.Z+1] != nil ||
				container.SmallBlocks[inner.X+1][inner.Y+1][inner.Z] != nil ||
				container.SmallBlocks[inner.X+1][inner.Y][inner.Z+1] != nil ||
				container.SmallBlocks[inner.X+1][inner.Y+1][inner.Z+1] != nil {
			return		
		}
	}

	container.SmallBlocks[inner.X][inner.Y][inner.Z] = bic
	
	if inner.X > 0 && container.SmallBlocks[inner.X-1][inner.Y][inner.Z] != nil && container.SmallBlocks[inner.X-1][inner.Y][inner.Z].Slope == SlopeNone {/*(
			container.SmallBlocks[inner.X-1][inner.Y][inner.Z].Slope == SlopeNone || 
			container.SmallBlocks[inner.X-1][inner.Y][inner.Z].Slope == SlopeNorthWest ||
			container.SmallBlocks[inner.X-1][inner.Y][inner.Z].Slope == SlopeSouthWest || 
			container.SmallBlocks[inner.X-1][inner.Y][inner.Z].Slope == SlopeWestTop ||
			container.SmallBlocks[inner.X-1][inner.Y][inner.Z].Slope == SlopeWestBottom) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionEast)
		
		if bic.Slope == SlopeNone {/*|| 
			bic.Slope == SlopeNorthEast ||
			bic.Slope == SlopeSouthEast || 
			bic.Slope == SlopeEastTop ||
			bic.Slope == SlopeEastBottom {*/
			container.SmallBlocks[inner.X-1][inner.Y][inner.Z].Adjacents = container.SmallBlocks[inner.X-1][inner.Y][inner.Z].Adjacents.Add(DirectionWest)
		} else {
			container.SmallBlocks[inner.X-1][inner.Y][inner.Z].Adjacents = container.SmallBlocks[inner.X-1][inner.Y][inner.Z].Adjacents.Remove(DirectionWest)
		}
	} 
	
	if inner.X < 3 && container.SmallBlocks[inner.X+1][inner.Y][inner.Z] != nil && container.SmallBlocks[inner.X+1][inner.Y][inner.Z].Slope == SlopeNone {/*(
			container.SmallBlocks[inner.X+1][inner.Y][inner.Z].Slope == SlopeNone || 
			container.SmallBlocks[inner.X+1][inner.Y][inner.Z].Slope == SlopeNorthEast ||
			container.SmallBlocks[inner.X+1][inner.Y][inner.Z].Slope == SlopeSouthEast || 
			container.SmallBlocks[inner.X+1][inner.Y][inner.Z].Slope == SlopeEastTop ||
			container.SmallBlocks[inner.X+1][inner.Y][inner.Z].Slope == SlopeEastBottom) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionWest)
		
		if bic.Slope == SlopeNone {/*|| 
			bic.Slope == SlopeNorthWest ||
			bic.Slope == SlopeSouthWest || 
			bic.Slope == SlopeWestTop ||
			bic.Slope == SlopeWestBottom {*/
			container.SmallBlocks[inner.X+1][inner.Y][inner.Z].Adjacents = container.SmallBlocks[inner.X+1][inner.Y][inner.Z].Adjacents.Add(DirectionEast)
		} else {
			container.SmallBlocks[inner.X+1][inner.Y][inner.Z].Adjacents = container.SmallBlocks[inner.X+1][inner.Y][inner.Z].Adjacents.Remove(DirectionEast)
		}
	} 
	
	if inner.Y > 0 && container.SmallBlocks[inner.X][inner.Y-1][inner.Z] != nil && container.SmallBlocks[inner.X][inner.Y-1][inner.Z].Slope == SlopeNone {/*(
			container.SmallBlocks[inner.X][inner.Y-1][inner.Z].Slope == SlopeNone || 
			container.SmallBlocks[inner.X][inner.Y-1][inner.Z].Slope == SlopeNorthTop ||
			container.SmallBlocks[inner.X][inner.Y-1][inner.Z].Slope == SlopeSouthTop || 
			container.SmallBlocks[inner.X][inner.Y-1][inner.Z].Slope == SlopeEastTop ||
			container.SmallBlocks[inner.X][inner.Y-1][inner.Z].Slope == SlopeWestTop) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionDown)
		
		if bic.Slope == SlopeNone {/*|| 
			bic.Slope == SlopeNorthBottom ||
			bic.Slope == SlopeSouthBottom || 
			bic.Slope == SlopeEastBottom ||
			bic.Slope == SlopeWestBottom {*/
			container.SmallBlocks[inner.X][inner.Y-1][inner.Z].Adjacents = container.SmallBlocks[inner.X][inner.Y-1][inner.Z].Adjacents.Add(DirectionUp)
		} else {
			container.SmallBlocks[inner.X][inner.Y-1][inner.Z].Adjacents = container.SmallBlocks[inner.X][inner.Y-1][inner.Z].Adjacents.Remove(DirectionUp)
		}
	} 
	
	if inner.Y < 3 && container.SmallBlocks[inner.X][inner.Y+1][inner.Z] != nil && container.SmallBlocks[inner.X][inner.Y+1][inner.Z].Slope == SlopeNone {/*(
			container.SmallBlocks[inner.X][inner.Y+1][inner.Z].Slope == SlopeNone || 
			container.SmallBlocks[inner.X][inner.Y+1][inner.Z].Slope == SlopeNorthBottom ||
			container.SmallBlocks[inner.X][inner.Y+1][inner.Z].Slope == SlopeSouthBottom || 
			container.SmallBlocks[inner.X][inner.Y+1][inner.Z].Slope == SlopeEastBottom ||
			container.SmallBlocks[inner.X][inner.Y+1][inner.Z].Slope == SlopeWestBottom) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionUp)
		
		if bic.Slope == SlopeNone {/*|| 
			bic.Slope == SlopeNorthTop ||
			bic.Slope == SlopeSouthTop || 
			bic.Slope == SlopeEastTop ||
			bic.Slope == SlopeWestTop {*/
			container.SmallBlocks[inner.X][inner.Y+1][inner.Z].Adjacents = container.SmallBlocks[inner.X][inner.Y+1][inner.Z].Adjacents.Add(DirectionDown)
		} else {
			container.SmallBlocks[inner.X][inner.Y+1][inner.Z].Adjacents = container.SmallBlocks[inner.X][inner.Y+1][inner.Z].Adjacents.Remove(DirectionDown)
		}
	} 
	
	if inner.Z > 0 && container.SmallBlocks[inner.X][inner.Y][inner.Z-1] != nil && container.SmallBlocks[inner.X][inner.Y][inner.Z-1].Slope == SlopeNone {/*(
			container.SmallBlocks[inner.X][inner.Y][inner.Z-1].Slope == SlopeNone || 
			container.SmallBlocks[inner.X][inner.Y][inner.Z-1].Slope == SlopeSouthBottom ||
			container.SmallBlocks[inner.X][inner.Y][inner.Z-1].Slope == SlopeSouthTop || 
			container.SmallBlocks[inner.X][inner.Y][inner.Z-1].Slope == SlopeSouthEast ||
			container.SmallBlocks[inner.X][inner.Y][inner.Z-1].Slope == SlopeSouthWest) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionSouth)
		
		if bic.Slope == SlopeNone {/* || 
			bic.Slope == SlopeNorthTop ||
			bic.Slope == SlopeNorthBottom || 
			bic.Slope == SlopeNorthEast ||
			bic.Slope == SlopeNorthWest {*/
			container.SmallBlocks[inner.X][inner.Y][inner.Z-1].Adjacents = container.SmallBlocks[inner.X][inner.Y][inner.Z-1].Adjacents.Add(DirectionNorth)
		} else {
			container.SmallBlocks[inner.X][inner.Y][inner.Z-1].Adjacents = container.SmallBlocks[inner.X][inner.Y][inner.Z-1].Adjacents.Remove(DirectionNorth)
		}
	} 
	
	if inner.Z < 3 && container.SmallBlocks[inner.X][inner.Y][inner.Z+1] != nil && container.SmallBlocks[inner.X][inner.Y][inner.Z+1].Slope == SlopeNone {/*(
			container.SmallBlocks[inner.X][inner.Y][inner.Z+1].Slope == SlopeNone || 
			container.SmallBlocks[inner.X][inner.Y][inner.Z+1].Slope == SlopeNorthBottom ||
			container.SmallBlocks[inner.X][inner.Y][inner.Z+1].Slope == SlopeNorthTop || 
			container.SmallBlocks[inner.X][inner.Y][inner.Z+1].Slope == SlopeNorthEast ||
			container.SmallBlocks[inner.X][inner.Y][inner.Z+1].Slope == SlopeNorthWest) {*/
		bic.Adjacents = bic.Adjacents.Add(DirectionNorth)
		
		if bic.Slope == SlopeNone {/*|| 
			bic.Slope == SlopeSouthTop ||
			bic.Slope == SlopeSouthBottom || 
			bic.Slope == SlopeSouthEast ||
			bic.Slope == SlopeSouthWest {*/
			container.SmallBlocks[inner.X][inner.Y][inner.Z+1].Adjacents = container.SmallBlocks[inner.X][inner.Y][inner.Z+1].Adjacents.Add(DirectionSouth)
		} else {
			container.SmallBlocks[inner.X][inner.Y][inner.Z+1].Adjacents = container.SmallBlocks[inner.X][inner.Y][inner.Z+1].Adjacents.Remove(DirectionSouth)
		}
	} 
			
	//this.addDirty(block.PackWithoutCoords())
	
	this.Unlock()
}

func (this *Chunk) RemoveBlockAt(coords WorldCoords) {
	log.Tracef("Chunk.RemoveBlockAt: %#v", coords)

	x,y,z,ok := this.CoordsToChunkCoords(coords)
	if !ok {
		return 
	}

	v := this.Blocks[x][y][z]
	if v != nil {
		this.Blocks[x][y][z] = nil
		//this.addDirty(v.PackWithoutCoords())
	}
}

type Chunks struct {
	sync.RWMutex
	Data map[PackedWorldCoords]*Chunk
}