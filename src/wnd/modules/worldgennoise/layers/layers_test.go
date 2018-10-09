package layers

import (
	"wnd/utils/log"
	"wnd/utils/test"

	"wnd/api"

	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"os/user"
	"testing"

	//"github.com/davecheney/profile"
)

const (
	sideSize int   = 32
	seed     int64 = 100500
)

func NoTestBiomeLayer(t *testing.T) {
	imgt := image.NewRGBA(image.Rect(0, 0, sideSize*sideSize, sideSize*sideSize))
	imgh := image.NewRGBA(image.Rect(0, 0, sideSize*sideSize, sideSize*sideSize))
	imgy := image.NewRGBA(image.Rect(0, 0, sideSize*sideSize, sideSize*sideSize))
	imgb := image.NewRGBA(image.Rect(0, 0, sideSize*sideSize, sideSize*sideSize))

	bl := NewBiomeLayer()

	for i := -(sideSize * int(api.ChunkSideSize)); i < (sideSize * int(api.ChunkSideSize)); i += int(api.ChunkSideSize) {
		for j := -(sideSize * int(api.ChunkSideSize)); j < (sideSize * int(api.ChunkSideSize)); j += int(api.ChunkSideSize) {
			for k := -(sideSize * int(api.ChunkSideSize)); k < (sideSize * int(api.ChunkSideSize)); k += int(api.ChunkSideSize) {
				chunk := &api.Chunk{Coords: api.WorldCoords{X: int32(i), Y: int32(j), Z: int32(k)}}

				//t.Logf("Biome layer generated: error %v , chunk: %#v", bl.Fill(seed, chunk), chunk)
				bl.Fill(seed, chunk)

				temp := uint8(chunk.BiomeData.Temperature + math.MaxInt8)
				hum := chunk.BiomeData.Humidity
				hei := uint8(chunk.BiomeData.SeaLevelHeight + math.MaxInt8)
				bias := chunk.BiomeData.HeightBias

				for b := 0; b < sideSize; b++ {
					for bb := 0; bb < sideSize; bb++ {
						imgt.SetRGBA(b+(i*sideSize/int(api.ChunkSideSize)), bb+(k*sideSize/int(api.ChunkSideSize)), color.RGBA{temp, 0, math.MaxUint8 - temp, math.MaxUint8})
						imgh.SetRGBA(b+(i*sideSize/int(api.ChunkSideSize)), bb+(k*sideSize/int(api.ChunkSideSize)), color.RGBA{hum, 0, math.MaxUint8 - hum, math.MaxUint8})
						imgy.SetRGBA(b+(i*sideSize/int(api.ChunkSideSize)), bb+(k*sideSize/int(api.ChunkSideSize)), color.RGBA{hei, 0, math.MaxUint8 - hei, math.MaxUint8})
						imgb.SetRGBA(b+(i*sideSize/int(api.ChunkSideSize)), bb+(k*sideSize/int(api.ChunkSideSize)), color.RGBA{bias, 0, math.MaxUint8 - bias, math.MaxUint8})
					}
				}
			}
		}
	}

	usr, _ := user.Current()

	filet, _ := os.Create(usr.HomeDir + "/biome-t.png")
	png.Encode(filet, imgt)
	filet.Close()

	fileh, _ := os.Create(usr.HomeDir + "/biome-h.png")
	png.Encode(fileh, imgh)
	fileh.Close()

	filey, _ := os.Create(usr.HomeDir + "/biome-y.png")
	png.Encode(filey, imgy)
	filey.Close()

	fileb, _ := os.Create(usr.HomeDir + "/biome-b.png")
	png.Encode(fileb, imgb)
	fileb.Close()
}

func TestTerrainLayer(t *testing.T) {
	//defer profile.Start(&profile.Config{CPUProfile:true, MemProfile:true, BlockProfile:true}).Stop()

	sideSize := 64
	offset := 0
	xOffset, zOffset := offset*int(api.ChunkSideSize), offset*int(api.ChunkSideSize)

	log.NewTestLogger(t)

	u := &api.Universe{}
	br := test.NewBlockRegistry(t)

	t.Error(br.Init())

	bl := NewBiomeLayer()
	tl := NewTerrainLayer(u, br.All())

	img := image.NewRGBA(image.Rect(0, 0, sideSize*int(api.ChunkSideSize), sideSize*int(api.ChunkSideSize)))

	for b := xOffset; b < xOffset+sideSize; b++ {
		for bb := zOffset; bb < zOffset+sideSize; bb++ {
			coords := api.WorldCoords{X: int32(b * int(api.ChunkSideSize)), Y: 64, Z: int32(bb * int(api.ChunkSideSize))}
			chunk := &api.Chunk{Coords: coords}

			bl.Fill(seed, chunk)
			tl.Fill(seed, chunk)

			min := 0.0
			max := 0.0

			for x := 0; x < int(api.ChunkSideSize); x++ {
				for z := 0; z < int(api.ChunkSideSize); z++ {
					for y := int(api.ChunkSideSize) - 1; y >= 0; y-- {
						crds := api.WorldCoords{X: coords.X + int32(x), Y: coords.Y + int32(y), Z: coords.Z + int32(z)}

						if b := chunk.BlockAt(crds); b != nil {
							max = math.Max(max, float64(y))
							min = math.Min(min, float64(y))
						}
					}
				}
			}

			max -= min

			for x := 0; x < int(api.ChunkSideSize); x++ {
				for z := 0; z < int(api.ChunkSideSize); z++ {
					var maxHeight int32 = math.MinInt32

					for y := int(api.ChunkSideSize) - 1; y >= 0; y-- {
						crds := api.WorldCoords{X: coords.X + int32(x), Y: coords.Y + int32(y), Z: coords.Z + int32(z)}

						if b := chunk.BlockAt(crds); b!= nil {
							maxHeight = crds.Y
							break
						}
					}

					h := uint8((float64(maxHeight) - min) / max * math.MaxUint8) //uint8((float64(maxHeight) + float64(math.MaxInt16)) / float64(math.MaxInt32) * float64(math.MaxUint8))

					//t.Logf("Test terrain %v: %d (%d) / rgba %d/%d", chunk.Coords(), h, maxHeight, int(coords.X) + x - (xOffset * int(api.ChunkSideSize)), int(coords.Z) + z - (zOffset * int(api.ChunkSideSize)))

					img.SetRGBA(int(coords.X)+x-(xOffset*int(api.ChunkSideSize)), int(coords.Z)+z-(zOffset*int(api.ChunkSideSize)), color.RGBA{h, h, h, math.MaxUint8})
				}
			}
		}
	}

	usr, _ := user.Current()

	file, _ := os.Create(usr.HomeDir + "/terrain.png")
	png.Encode(file, img)
	file.Close()
}
