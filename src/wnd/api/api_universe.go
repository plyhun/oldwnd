package api

import (
	"wnd/utils/log"
)

const (
	dayLength uint64 = 1000 * 60 * 10 
)

type Universe struct {
	Chunks
	
	Name string
	Seed int64
	Size uint32
	Age  uint64

	Entities map[string]*Entity
}

/*func (this *Universe) IsCoordAtTheEdge(c Coords, visibleChunks uint32) bool {
	edge := float32((this.Size - visibleChunks + 1) * uint32(ChunkSideSize))
	
	return c.X >= edge || c.X <= -edge || c.Z >= edge || c.Z <= -edge || c.Y >= edge || c.Y <= -edge 
}

func (this *Universe) IsWorldCoordAtTheEdge(c WorldCoords, visibleChunks uint32) bool {
	edge := int32((this.Size - visibleChunks + 1) * uint32(ChunkSideSize))
	
	return c.X >= edge || c.X <= -edge || c.Z >= edge || c.Z <= -edge || c.Y >= edge || c.Y <= -edge 
}

func (this *Universe) GetWorldCoordsForTiling(c WorldCoords, visibleChunks uint32) map[AdjacencyList]WorldCoords {
	edge := int32((this.Size - visibleChunks + 1) * uint32(ChunkSideSize))
	
	result := make(map[AdjacencyList]WorldCoords, 7)
	
	delta := int32(2 * ChunkSideSize * this.Size)
	
	dir := AdjacencyList(0)
	deltaCoords := WorldCoords{0,0,0}
	
	if c.X >= edge {
		//result[AdjacencyList(0).Add(DirectionEast)] = WorldCoords{X: c.X - (2 * ChunkSideSize * int32(this.Size))}
		
		dir = dir.Add(DirectionEast)
		deltaCoords = WorldCoords{X: -delta, Y: 0, Z: 0}
	} else if c.X <= -edge {
		//result[AdjacencyList(0).Add(DirectionWest)] = WorldCoords{X: c.X + (2 * ChunkSideSize * int32(this.Size))}
		
		dir = dir.Add(DirectionWest)
		deltaCoords = WorldCoords{X: delta, Y: 0, Z: 0}
	}
	
	if dir.IsDefined() {
		for k,v := range result {
			dir1 := dir.Add(k.Parse()...)
			coords1 := WorldCoords{X: deltaCoords.X + v.X, Y: deltaCoords.Y + v.Y, Z: deltaCoords.Z + v.Z}
			
			result[dir1] = coords1
		}
		
		result[dir] = WorldCoords{X: c.X + deltaCoords.X, Y: c.Y + deltaCoords.Y, Z: c.Z + deltaCoords.Z}
	}
	
	dir = AdjacencyList(0)
	deltaCoords = WorldCoords{0,0,0}
	
	if c.Z >= edge {
		//result[AdjacencyList(0).Add(DirectionSouth)] = WorldCoords{X: c.Z - (2 * ChunkSideSize * int32(this.Size))}
		
		dir = dir.Add(DirectionSouth)
		deltaCoords = WorldCoords{X: 0, Y: 0, Z: -delta}
	} else if c.Z <= -edge {
		//result[AdjacencyList(0).Add(DirectionNorth)] = WorldCoords{X: c.Z + (2 * ChunkSideSize * int32(this.Size))}
		
		dir = dir.Add(DirectionNorth)
		deltaCoords = WorldCoords{X: 0, Y: 0, Z: delta}
	}
	
	if dir.IsDefined() {
		for k,v := range result {
			dir1 := dir.Add(k.Parse()...)
			coords1 := WorldCoords{X: deltaCoords.X + v.X, Y: deltaCoords.Y + v.Y, Z: deltaCoords.Z + v.Z}
			
			result[dir1] = coords1
		}
		
		result[dir] = WorldCoords{X: c.X + deltaCoords.X, Y: c.Y + deltaCoords.Y, Z: c.Z + deltaCoords.Z}
	}
	
	dir = AdjacencyList(0)
	deltaCoords = WorldCoords{0,0,0}
	
	if c.Y >= edge {
		//result[AdjacencyList(0).Add(DirectionDown)] = WorldCoords{X: c.Y - (2 * ChunkSideSize * int32(this.Size))}
		
		dir = dir.Add(DirectionDown)
		deltaCoords = WorldCoords{X: 0, Y: -delta, Z: 0}
	} else if c.Y <= -edge {
		//result[AdjacencyList(0).Add(DirectionUp)] = WorldCoords{X: c.Y + (2 * ChunkSideSize * int32(this.Size))}
		
		dir = dir.Add(DirectionUp)
		deltaCoords = WorldCoords{X: 0, Y: delta, Z: 0}
	}
	
	if dir.IsDefined() {
		for k,v := range result {
			dir1 := dir.Add(k.Parse()...)
			coords1 := WorldCoords{X: deltaCoords.X + v.X, Y: deltaCoords.Y + v.Y, Z: deltaCoords.Z + v.Z}
			
			result[dir1] = coords1
		}
		
		result[dir] = WorldCoords{X: c.X + deltaCoords.X, Y: c.Y + deltaCoords.Y, Z: c.Z + deltaCoords.Z}
	}
	
	log.Warnf("%v -> %v", c, result)
	
	return result
}

func (this *Universe) GetRightCoords(in Coords) Coords {
	log.Warnf("in %v", in)
	
	if in.X > float32(this.Size) {
		diff := in.X - float32(this.Size)
		in.X = diff - float32(this.Size)
	} else if in.X < -float32(this.Size) {
		diff := in.X + float32(this.Size)
		in.X = float32(this.Size) + diff
	}

	if in.Z > float32(this.Size) {
		diff := in.Z - float32(this.Size)
		in.Z = diff - float32(this.Size)
	} else if in.Z < -float32(this.Size) {
		diff := in.Z + float32(this.Size)
		in.Z = float32(this.Size) + diff
	}

	if in.Y > float32(this.Size) {
		diff := in.Y - float32(this.Size)
		in.Y = diff - float32(this.Size)
	} else if in.Y < -float32(this.Size) {
		diff := in.Y + float32(this.Size)
		in.Y = float32(this.Size) + diff
	}
	
	log.Warnf("out %v", in)

	return in
}

func (this *Universe) GetRightWorldCoords(in WorldCoords) WorldCoords {
	log.Warnf("in %v", in)
	
	if int(in.X) > int(this.Size) {
		diff := int(in.X) - int(this.Size)
		in.X = int32(diff - int(this.Size))
	} else if int(in.X) < -int(this.Size) {
		diff := int(in.X) + int(this.Size)
		in.X = int32(int(this.Size) + diff)
	}

	if int(in.Z) > int(this.Size) {
		diff := int(in.Z) - int(this.Size)
		in.Z = int32(diff - int(this.Size))
	} else if int(in.Z) < -int(this.Size) {
		diff := int(in.Z) + int(this.Size)
		in.Z = int32(int(this.Size) + diff)
	}

	if int(in.Y) > int(this.Size) {
		diff := int(in.Y) - int(this.Size)
		in.Y = int32(diff - int(this.Size))
	} else if int(in.Y) < -int(this.Size) {
		diff := int(in.Y) + int(this.Size)
		in.Y = int32(int(this.Size) + diff)
	}
	
	log.Warnf("out %v", in)

	return in
}*/

func (this *Universe) ID() string {
	return this.Name
}

func (this *Universe) AddBlock(coords WorldCoords, b *Block) *Chunk {
	c := this.Chunks.Data[coords.ChunkCoords().Pack()]
	if c != nil {
		c.SetBlock(coords, b)
		return c
	}
	
	return nil
}

func (this *Universe) BlockAt(coords WorldCoords) *Block {
	chunkCoords := coords.ChunkCoords()
	
	if c,ok := this.Chunks.Data[chunkCoords.Pack()]; ok {
		b := c.BlockAt(coords)
		if b != nil {
			return b.Block
		} else {
			return nil
		}
	} else {
		return nil
	}
}

type World struct {
	*Universe
	
	Usage string
	Owner string
	
	VisibleRadius uint32
}

func (this *World) Target() string {
	return this.Usage
}

func (this *World) GetShift(target,viewer Coords) (ret Coords, ok bool) {
	ret = Coords{0,0,0}
	
	edge := float32((this.Size - this.VisibleRadius) * uint32(ChunkSideSize))
	
	if (viewer.X <= edge && viewer.X > -edge) && (viewer.Z <= edge && viewer.Z > -edge) {
		ok = false
		return
	}
	
	if (target.X <= edge && target.X > -edge) && (target.Z <= edge && target.Z > -edge) {
		ok = false
		return
	} 
	
	blockSize := float32(this.Size * 2 * uint32(ChunkSideSize))
	
	if (target.X > 0 && viewer.X < 0) {
		ret.X -= blockSize
	} else if (target.X < 0 && viewer.X > 0) {
		ret.X += blockSize
	}
	
	if (target.Z > 0 && viewer.Z < 0) {
		ret.Z -= blockSize
	} else if (target.Z < 0 && viewer.Z > 0) {
		ret.Z += blockSize
	}
	
	shifted := Coords{X: target.X+ret.X, Y: target.Y+ret.Y, Z: target.Z+ret.Z}
	blockRadius := float32((this.VisibleRadius + 1) * uint32(ChunkSideSize))
	
	if viewer.X - shifted.X > blockRadius || viewer.X - shifted.X < -blockRadius || viewer.Z - shifted.Z > blockRadius || viewer.Z - shifted.Z < -blockRadius {
		ok = false
	} else if !shifted.Equals(target){
		ok = true
		log.Debugf("shift %v -> %v for %v", target, ret, viewer)
	}
	
	return
}
