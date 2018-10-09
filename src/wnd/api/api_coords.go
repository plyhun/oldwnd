package api

import (
	"fmt"
)

const worldCoordsPackFormat string = "%v/%v/%v"
const WorldCoordSize float32 = 1.0

type WorldCoords struct {
	X, Y, Z int32
}

type Vector struct {
	X, Y, Z float32
}

type Coords struct {
	X, Y, Z float32
}

type PackedWorldCoords string

func PackWorldCoords(x, y, z int32) PackedWorldCoords {
	return PackedWorldCoords(fmt.Sprintf(worldCoordsPackFormat, x, y, z))
}

func (this WorldCoords) Pack() PackedWorldCoords {
	return PackWorldCoords(this.X, this.Y, this.Z)
}

func UnpackWorldCoords(pc PackedWorldCoords) (coords WorldCoords, e error) {
	coords = WorldCoords{}
	_, e = fmt.Sscanf(string(pc), worldCoordsPackFormat, &coords.X, &coords.Y, &coords.Z)

	return coords, e
}

func (this WorldCoords) ToCoords() Coords {
	return Coords{X: float32(this.X) * WorldCoordSize, Y: float32(this.Y) * WorldCoordSize, Z: float32(this.Z) * WorldCoordSize}
}

func (this Coords) ToWorldCoords() WorldCoords {
	return WorldCoords{X: int32(this.X / WorldCoordSize), Y: int32(this.Y / WorldCoordSize), Z: int32(this.Z / WorldCoordSize)}
}

func (this WorldCoords) Equals(other WorldCoords) bool {
	return this.X == other.X && this.Y == other.Y && this.Z == other.Z
}

func (this Coords) Equals(other Coords) bool {
	return this.X == other.X && this.Y == other.Y && this.Z == other.Z
}

func (this WorldCoords) GetNeighbor(directions ...Direction) WorldCoords {
	res := WorldCoords{X: this.X, Y: this.Y, Z: this.Z}

	for _, d := range directions {
		v := d.Vector()
		res.X += v.X
		res.Y += v.Y
		res.Z += v.Z
	}

	return res
}
