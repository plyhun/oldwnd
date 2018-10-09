package layers

import (
	"wnd/api"
	"wnd/utils/log"
	
	"math"
)

type BiomeLayer struct {
	temperature perlinGenerator
	humidity perlinGenerator
	height perlinGenerator
	heightbias perlinGenerator
}

func NewBiomeLayer() *BiomeLayer {
	return &BiomeLayer{
		temperature: perlinGenerator{
			octaves: func(seed int) int {
				return 1
			},
			scale: func(seed float64) float64 {
				return 8 * api.ChunkSideSizeFloat64
			},
			seeder: func(seed int64) int64 {
				return (seed / 5) * 3 
			},
			alpha: func(seed float64) float64 {
				return 2
			},
			beta: func(seed float64) float64 {
				return 2
			},
		},
		humidity: perlinGenerator{
			octaves: func(seed int) int {
				return 1
			},
			scale: func(seed float64) float64 {
				return 8 * api.ChunkSideSizeFloat64
			},
			seeder: func(seed int64) int64 {
				return seed / 4 
			},
			alpha: func(seed float64) float64 {
				return 3
			},
			beta: func(seed float64) float64 {
				return 6
			},
		},
		height: perlinGenerator{
			octaves: func(seed int) int {
				return 1
			},
			scale: func(seed float64) float64 {
				return 8 * api.ChunkSideSizeFloat64
			},
			seeder: func(seed int64) int64 {
				return seed / 2
			},
			alpha: func(seed float64) float64 {
				return 84
			},
			beta: func(seed float64) float64 {
				return 0.11
			},
		},
		heightbias: perlinGenerator{
			octaves: func(seed int) int {
				return 1
			},
			scale: func(seed float64) float64 {
				return 10 * api.ChunkSideSizeFloat64
			},
			seeder: func(seed int64) int64 {
				return seed / 3
			},
			alpha: func(seed float64) float64 {
				return 0.4
			},
			beta: func(seed float64) float64 {
				return 0.31
			},
		},
	}
}

func (this *BiomeLayer) ID() string {
	return "BiomeLayer"
}

func (this *BiomeLayer) Fill(seed int64, chunk *api.Chunk, size uint32) (e error) {
	t := this.temperature.fill(chunk.Coords.X, chunk.Coords.Z, seed, 0,0,0,0)
	h := this.humidity.fill(chunk.Coords.X, chunk.Coords.Z, seed, 0,0,0,0)
	y := this.height.fill(chunk.Coords.X, chunk.Coords.Z, seed, 0,0,0,0)
	b := this.heightbias.fill(chunk.Coords.X, chunk.Coords.Z, seed, 0,0,0,0)
	
	log.Debugf("t %v h %v y %v b %v", t, h, y, b)
	
	chunk.BiomeData.Temperature = int8((t * math.MaxFloat64) / (math.MaxFloat64 / math.MaxInt8))
	chunk.BiomeData.Humidity = uint8((h * math.MaxFloat64) / (math.MaxFloat64 / math.MaxInt8) + math.MaxInt8)
	chunk.BiomeData.SeaLevelHeight = int8((y * math.MaxFloat64) / (math.MaxFloat64 / math.MaxInt8))
	chunk.BiomeData.HeightBias = uint8((b * math.MaxFloat64) / (math.MaxFloat64 / math.MaxInt8) + math.MaxInt8)
	
	log.Debugf("BiomeLayer.fill: seed %d chunk %#v", seed, chunk.BiomeData)
	
	return
}