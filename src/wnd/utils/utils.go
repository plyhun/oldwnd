package utils

import (
	"wnd/api"
	"wnd/utils/log"

	"os"
	"path/filepath"
	"sort"
)

func GetAppDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Errorf("Cannot detect own application dir: %v", err)
		return ""
	}

	return dir
}

func CheckAndMakeDir(dir string) error {
	log.Tracef("%s", dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0777)
	} else {
		return nil
	}
}

func GetRightWorldCoords(in api.WorldCoords, size uint32) (api.WorldCoords, bool) {
	log.Tracef("in %v", in)
	
	processed := false

	size *= api.ChunkSideSize

	if int(in.X) >= int(size) {
		processed = true
		diff := int(in.X) - int(size)
		in.X = int32(diff - int(size))
	} else if int(in.X) < -int(size) {
		processed = true
		diff := int(in.X) + int(size)
		in.X = int32(int(size) + diff)
	}

	if int(in.Z) >= int(size) {
		processed = true
		diff := int(in.Z) - int(size)
		in.Z = int32(diff - int(size))
	} else if int(in.Z) < -int(size) {
		processed = true
		diff := int(in.Z) + int(size)
		in.Z = int32(int(size) + diff)
	}

	/*if int(in.Y) >= int(size) {
		processed = true
		diff := int(in.Y) - int(size)
		in.Y = int32(diff - int(size))
	} else if int(in.Y) < -int(size) {
		processed = true
		diff := int(in.Y) + int(size)
		in.Y = int32(int(size) + diff)
	}*/

	log.Tracef("out %v", in)

	return in, processed
}

func GetRightCoords(in api.Coords, size uint32) (api.Coords, bool) {
	log.Tracef("in %v", in)
	
	processed := false

	size *= uint32(int32(api.ChunkSideSize))

	if in.X >= float32(size) {
		processed = true
		diff := in.X - float32(size)
		in.X = diff - float32(size)
	} else if in.X < -float32(size) {
		processed = true
		diff := in.X + float32(size)
		in.X = float32(size) + diff
	}

	if in.Z >= float32(size) {
		processed = true
		diff := in.Z - float32(size)
		in.Z = diff - float32(size)
	} else if in.Z < -float32(size) {
		processed = true
		diff := in.Z + float32(size)
		in.Z = float32(size) + diff
	}

	/*if in.Y >= float32(size) {
		processed = true
		diff := in.Y - float32(size)
		in.Y = diff - float32(size)
	} else if in.Y < -float32(size) {
		processed = true
		diff := in.Y + float32(size)
		in.Y = float32(size) + diff
	}*/

	log.Tracef("out %v", in)

	return in, processed
}

func EntityGetRequiredAndPurgedChunks(oldChunkCoords, newChunkCoords api.WorldCoords, radius uint8) ([]api.WorldCoords, []api.WorldCoords) {
	log.Tracef("old position %v, new position %v", oldChunkCoords, newChunkCoords)

	sideSize := 1 + (2 * int(radius))

	var oldChunks, newChunks []api.WorldCoords

	if !oldChunkCoords.Equals(newChunkCoords) {

		xshift, yshift, zshift := newChunkCoords.X-oldChunkCoords.X, newChunkCoords.Y-oldChunkCoords.Y, newChunkCoords.Z-oldChunkCoords.Z

		switchedChunksCount := 0
		sqr := int(sideSize * sideSize)

		if xshift != 0 {
			switchedChunksCount += sqr
		}

		if yshift != 0 {
			switchedChunksCount += sqr
		}

		if zshift != 0 {
			switchedChunksCount += sqr
		}

		switch {
		case switchedChunksCount >= sqr*3:
			switchedChunksCount -= (int(sideSize) * 3) - 1
		case switchedChunksCount >= sqr*2:
			switchedChunksCount -= int(sideSize)
		}

		oldMinMaxChunks := []api.WorldCoords{
			api.WorldCoords{X: oldChunkCoords.X - int32(radius)*int32(api.ChunkSideSize), Y: oldChunkCoords.Y - int32(radius)*int32(api.ChunkSideSize), Z: oldChunkCoords.Z - int32(radius)*int32(api.ChunkSideSize)},
			api.WorldCoords{X: oldChunkCoords.X + int32(radius+1)*int32(api.ChunkSideSize), Y: oldChunkCoords.Y + int32(radius+1)*int32(api.ChunkSideSize), Z: oldChunkCoords.Z + int32(radius+1)*int32(api.ChunkSideSize)},
		}

		newMinMaxChunks := []api.WorldCoords{
			api.WorldCoords{X: newChunkCoords.X - int32(radius)*int32(api.ChunkSideSize), Y: newChunkCoords.Y - int32(radius)*int32(api.ChunkSideSize), Z: newChunkCoords.Z - int32(radius)*int32(api.ChunkSideSize)},
			api.WorldCoords{X: newChunkCoords.X + int32(radius+1)*int32(api.ChunkSideSize), Y: newChunkCoords.Y + int32(radius+1)*int32(api.ChunkSideSize), Z: newChunkCoords.Z + int32(radius+1)*int32(api.ChunkSideSize)},
		}

		//log.Warnf("old chunk %v (%v), new chunk %v (%v) => %v", oldChunkCoords, oldMinMaxChunks, newChunkCoords, newMinMaxChunks, switchedChunksCount)

		oldChunks, newChunks = make([]api.WorldCoords, switchedChunksCount), make([]api.WorldCoords, switchedChunksCount)
		index := 0

		for x := oldMinMaxChunks[0].X; x < oldMinMaxChunks[1].X; x += int32(api.ChunkSideSize) {
			for y := oldMinMaxChunks[0].Y; y < oldMinMaxChunks[1].Y; y += int32(api.ChunkSideSize) {
				for z := oldMinMaxChunks[0].Z; z < oldMinMaxChunks[1].Z; z += int32(api.ChunkSideSize) {
					if x < newMinMaxChunks[0].X || x >= newMinMaxChunks[1].X ||
						y < newMinMaxChunks[0].Y || y >= newMinMaxChunks[1].Y ||
						z < newMinMaxChunks[0].Z || z >= newMinMaxChunks[1].Z {
						oldChunks[index] = api.WorldCoords{X: x, Y: y, Z: z}

						delta := api.WorldCoords{X: oldChunkCoords.X - x, Y: oldChunkCoords.Y - y, Z: oldChunkCoords.Z - z}

						newChunks[index] = api.WorldCoords{X: delta.X + newChunkCoords.X, Y: delta.Y + newChunkCoords.Y, Z: delta.Z + newChunkCoords.Z}

						//log.Warnf("!!! %v -> %v", oldChunks[index], newChunks[index])
						index++
					}
				}
			}
		}
	} else {
		oldChunks, newChunks = make([]api.WorldCoords, 0), make([]api.WorldCoords, 0)
	}

	if len(oldChunks) > 0 || len(newChunks) > 0 {
		//log.Warnf("total: old %v, new %v", oldChunks, newChunks)
	}

	sort.Sort(sortableCoordsArray{data: newChunks, base: newChunkCoords})

	return oldChunks, newChunks
}

func Int32Min(x, y int32) int32 {
	if x < y {
		return x
	} else {
		return y
	}
}

func Int32Max(x, y int32) int32 {
	if x > y {
		return x
	} else {
		return y
	}
}

type sortableCoordsArray struct {
	data []api.WorldCoords
	base api.WorldCoords
}

func (this sortableCoordsArray) Len() int {
	return len(this.data)
}
func (this sortableCoordsArray) Less(i, j int) bool {
	di, dj := api.WorldCoords{X: this.data[i].X - this.base.X, Y: this.data[i].Y - this.base.Y, Z: this.data[i].Z - this.base.Z}, api.WorldCoords{X: this.data[j].X - this.base.X, Y: this.data[j].Y - this.base.Y, Z: this.data[j].Z - this.base.Z}

	if di.X < 0 {
		di.X = -di.X
	}
	if di.Y < 0 {
		di.Y = -di.Y
	}
	if di.Z < 0 {
		di.Z = -di.Z
	}

	if dj.X < 0 {
		dj.X = -dj.X
	}
	if dj.Y < 0 {
		dj.Y = -dj.Y
	}
	if dj.Z < 0 {
		dj.Z = -dj.Z
	}

	return di.X <= dj.X && di.Y <= dj.Y && di.Z <= dj.Z
}
func (this sortableCoordsArray) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}
