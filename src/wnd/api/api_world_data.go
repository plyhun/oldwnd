package api

import "fmt"

const ChunkSideSizePower = 4
const ChunkSideSize uint32 = uint32(1 << ChunkSideSizePower)
const ChunkLastBlock int32 = int32(ChunkSideSize - 1)
const ChunkMiddleBlock int32 = int32(ChunkSideSize >> 2)
const ChunkSideSizeFloat64 = float64(ChunkSideSize)

type BlockSize uint8
type BlockSlope uint8
type BlockOrientation uint8
type Direction uint8

const (
	SizeFull BlockSize = iota
	SizeHalf
	SizeQuarter
)

const (
	SlopeNone BlockSlope = iota
	SlopeNorthTop
	SlopeNorthBottom
	SlopeSouthTop
	SlopeSouthBottom
	SlopeEastTop
	SlopeEastBottom
	SlopeWestTop
	SlopeWestBottom
	SlopeNorthEast
	SlopeNorthWest
	SlopeSouthEast
	SlopeSouthWest
)

const (
	OrientationNorthTop BlockOrientation = iota
	OrientationNorthBottom
	OrientationSouthTop
	OrientationSouthBottom
	OrientationEastTop
	OrientationEastBottom
	OrientationWestTop
	OrientationWestBottom
	OrientationNorthEast
	OrientationNorthWest
	OrientationSouthEast
	OrientationSouthWest
	OrientationDefault = OrientationNorthTop
)

const (
	DirectionNorth Direction = iota
	DirectionSouth
	DirectionWest
	DirectionEast
	DirectionUp
	DirectionDown
)

var DirectionAll = [6]Direction{
	DirectionNorth,
	DirectionSouth,
	DirectionWest,
	DirectionEast,
	DirectionUp,
	DirectionDown,
}

func (this BlockSize) Multiplier() float32 {
	switch this {
	case SizeQuarter:
		return WorldCoordSize / 4
	case SizeHalf:
		return WorldCoordSize / 2
	default:
		return WorldCoordSize
	}
}

func (this Direction) Vector() WorldCoords {
	res := WorldCoords{X: 0, Y: 0, Z: 0}
	switch this {
	case DirectionNorth:
		res.Z++
	case DirectionSouth:
		res.Z--
	case DirectionUp:
		res.Y++
	case DirectionDown:
		res.Y--
	case DirectionEast:
		res.X--
	case DirectionWest:
		res.X++
	}
	return res
}

type BiomeData struct {
	Temperature    int8
	Humidity       uint8
	HeightBias     uint8
	SeaLevelHeight int8
}

type BlockDefinition struct {
	ID                    uint32
	Type                  string
	Name                  LocalizableString
	Sizes                 []BlockSize
	Transparency          uint8
	Hardness              uint8
	Fluidity              uint8
	EntityHealthInfluence int16
	Slipperiness          int16
	Mass                  int32
	LightEmission         uint32
	Orientable            bool
	Slopeable             bool
}

type Tileable struct {
	Adjacents AdjacencyList
}

type Block struct {
	ID          uint32
	Size        BlockSize
	Slope       BlockSlope
	Orientation BlockOrientation
	Metadata    map[string]interface{}
}

type PackedBlock string

const blockDataPackFormat string = "%d %d %d %d"

func (this *Block) PackWithoutCoords() PackedBlock {
	return PackedBlock(fmt.Sprintf(blockDataPackFormat, this.ID, this.Size, this.Slope, this.Orientation))
}

func UnpackWithoutCoords(pbd PackedBlock) (bd *Block, e error) {
	bd = new(Block)
	_, e = fmt.Sscanf(string(pbd), blockDataPackFormat, &bd.ID, &bd.Size, &bd.Slope, &bd.Orientation)
	return bd, e
}

func (this WorldCoords) ChunkCoords() WorldCoords {
	chunkCoords := this

	diff := chunkCoords.X % int32(ChunkSideSize)
	chunkCoords.X -= int32(diff)
	if diff < 0 {
		chunkCoords.X -= int32(ChunkSideSize)
	}

	diff = chunkCoords.Y % int32(ChunkSideSize)
	chunkCoords.Y -= int32(diff)
	if diff < 0 {
		chunkCoords.Y -= int32(ChunkSideSize)
	}
	
	diff = chunkCoords.Z % int32(ChunkSideSize)
	chunkCoords.Z -= int32(diff)
	if diff < 0 {
		chunkCoords.Z -= int32(ChunkSideSize)
	}
	
	//log.Tracef("%v => %v", this, chunkCoords)

	return chunkCoords
}
