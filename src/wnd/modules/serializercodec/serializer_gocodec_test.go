package serializercodec

import (
	"wnd/api"
	"wnd/modules"
	"wnd/api/events"

	"testing"
	
	"github.com/kr/pretty"
)

//var lizer Serializer = serealimpl.NewSerealSerializer()
var lizer modules.Serializer = New()
var (
	testBlockData = &api.BlockData{
		ID:         8,
		Size:        api.SizeFull,
		Slope:       api.SlopeNone,
		Orientation: api.OrientationDefault,
	}
	
	testBlock = &api.Block{
		Type:                  "ablock",
		Name:                  api.LocalizableString("a block"),
		Sizes:                 []api.BlockSize{api.SizeFull, api.SizeHalf},
		Transparency:          0,
		Hardness:              255,
		Fluidity:              0,
		EntityHealthInfluence: 0,
		Slipperiness:          0,
		Mass:                  11,
		LightEmission:         0,
		Orientable:            true,
		Slopeable:             false,
		ID:                    0xcafebabe,
	}
	
	atestChunk = &api.Chunk{
		Coords: api.WorldCoords{X: 16, Y: 32, Z: 64},
		BiomeData: api.BiomeData{
			Temperature: 0,
			Humidity: 10,
			HeightBias: 20,
			SeaLevelHeight: 30,
		},
	}
	
	testChunksEvent = &events.Chunk{
		General: events.General{
			EventID: "event-chunks",
			EventTime: 0xcafebabe,
			EventUniverseID: "theUniverse",
		},
		Chunk: atestChunk,
		Coords: atestChunk.Coords,
	}
)

func init() {
	for x := 0; x < int(api.ChunkSideSize); x++ {
		for y := 0; y < int(api.ChunkSideSize); y++ {
			for z := 0; z < int(api.ChunkSideSize); z++ {
				//if (x + y + z) % 16 > 0 {
					atestChunk.Blocks[x][y][z] = &api.BlockData{
						ID:         uint32(x+y+z),
						Size:        api.SizeFull,
						Slope:       api.SlopeNone,
						Orientation: api.OrientationDefault,
					}
				//}
			}
		}
	}
}

func aTestBlockSerializer(t *testing.T) {
	t.Logf("\n\n\n\n>>>>>>>Test block serializer: %v", lizer.Init())
	
	b,e := lizer.Serialize(testBlock)

	t.Logf("Data: %+v\n error %v", b, e)

	mm,ee := lizer.Deserialize(b)

	t.Logf("Got: %# v\n error %v", pretty.Formatter(mm), ee)

	lizer.Destroy()
	t.Logf("\n\n\n\n<<<<<<<<")
}

func TestChunkSerializer(t *testing.T) {
	t.Logf("\n\n\n\n>>>>>>>Test chunk serializer: %v", lizer.Init())
	
	b,e := lizer.Serialize(atestChunk)
	
	t.Logf("Data: %+v\n error %v", b, e)

	mm,ee := lizer.Deserialize(b)

	t.Logf("Got: %+v\n error %v", pretty.Formatter(mm), ee)
	
	lizer.Destroy()
	t.Logf("\n\n\n\n<<<<<<<<")
}

func aTestChunksEventSerializer(t *testing.T) {
	t.Logf("\n\n\n\n>>>>>>>Test chunk event serializer: %v", lizer.Init())
	
	b,e := lizer.Serialize(testChunksEvent)
	
	t.Logf("Data: data %+v\n  error %v len %v ", b, e, len(b))

	mm,ee := lizer.Deserialize(b)

	t.Logf("Got: %#v\n error %v", mm, ee)
	
	lizer.Destroy()
	t.Logf("\n\n\n\n<<<<<<<<")
}