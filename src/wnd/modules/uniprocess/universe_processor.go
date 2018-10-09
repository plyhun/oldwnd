package uniprocess

import (
	"wnd/api"
	"wnd/api/events"
	"wnd/base/entity"
	"wnd/modules"
	"wnd/utils"
	"wnd/utils/log"

	"reflect"
	"sync"
	"time"
	
	"github.com/karlseguin/ccache"
)

const (
	blockDataLengthForEvent = int(int32(api.ChunkSideSize))
)

func New() modules.UniverseProcessor {
	return &uniProcessorImpl{
		deferredProcessResults: make(chan api.Outputable),
		players:                make(map[string]*entity.Playable),
		lru: ccache.New(ccache.Configure()),
	}
}

type uniProcessorImpl struct {
	uniLock sync.Mutex

	deferredProcessResults chan api.Outputable
	players                map[string]*entity.Playable

	Universe       *api.Universe          `inject:""`
	WorldGenerator modules.WorldGenerator `inject:""`
	BlockRegistry  modules.BlockRegistry  `inject:""`
	
	lru *ccache.Cache
}

func (this *uniProcessorImpl) ID() string {
	return "universeProcessor"
}

func (this *uniProcessorImpl) Priority() int8 {
	return 10
}

func (this *uniProcessorImpl) Chunks(coords []api.WorldCoords) []*api.Chunk {
	chunks := make([]*api.Chunk, len(coords))

	for i, coord := range coords {
		chunks[i] = this.getOrCreateChunk(coord)
	}

	return chunks
}

func (this *uniProcessorImpl) Blocks(coords []api.WorldCoords) []*api.Block {
	blocks := make([]*api.Block, len(coords))

	for i, coord := range coords {
		c := this.getOrCreateChunk(coord.ChunkCoords())

		b := c.BlockAt(coord)

		if b == nil {
			blocks[i] = nil
		} else {
			blocks[i] = b.Block
		}
	}

	return blocks
}

func (this *uniProcessorImpl) World(entityId string) *api.World {
	log.Tracef("entity id# %s", entityId)

	e := this.Entity(entityId)

	if e == nil {
		log.Warnf("no entity found with ID# %v", entityId)
		return nil
	}

	this.players[entityId] = entity.NewPlayable(e)

	w := &api.World{
		Owner: entityId,
		Universe: &api.Universe{
			Seed:   this.Universe.Seed,
			Size:   this.Universe.Size,
			Age:    this.Universe.Age,
			Name:   this.Universe.Name,
			Chunks: api.Chunks{Data: make(map[api.PackedWorldCoords]*api.Chunk)},
		},
	}

	/*for _,c := range this.entityChunks(entityId) {
		w.Chunks[c.Coords.Pack()] = c
	}*/

	return w
}

func (this *uniProcessorImpl) Entity(entityId string) *api.Entity {
	log.Tracef("id# %s (all %v)", entityId, this.Universe)

	return this.Universe.Entities[entityId]
}

func (this *uniProcessorImpl) AddEntity(entity *api.Entity) {
	if entity == nil || entity.ID == "" {
		log.Errorf("invalid entity to add: %# v", entity)
		return
	}

	this.Universe.Entities[entity.ID] = entity
}

func (this *uniProcessorImpl) getEntityActualOffset(entityId string, offset api.Coords) (newOffset api.Coords) {

	newOffset = offset

	//TODO make collision check better
	/*e := this.Entity(entityId)

	if e == nil {
		log.Errorf("no entity to move: %#v", entityId)
		return
	}

	newCoords := api.Coords{X: e.Position.X + offset.X, Y: e.Position.Y + offset.Y, Z: e.Position.Z + offset.Z}
	newWorldCoords := newCoords.ToWorldCoords()

	if len(this.Blocks([]api.WorldCoords{newWorldCoords})) > 0 {
		newOffset = api.Coords{0, 0, 0}
	}*/

	return
}

func (this *uniProcessorImpl) entityChunks(entityId string, source string, extraChunksToLoadForEntity int32, metadata interface{}) []*api.Chunk {
	log.Tracef("%s", entityId)

	//entityWorldCoordsRadius := float32(extraChunksToLoadForEntity * int(int32(api.ChunkSideSize)))

	chunksStorageLength := int(1 + (2 * extraChunksToLoadForEntity))
	worldChunks := make([]*api.Chunk, chunksStorageLength*chunksStorageLength*chunksStorageLength)
	e := this.Entity(entityId)

	position := e.Position.ToWorldCoords().ChunkCoords()
	coords := position

	index, amountHoriz := 0, 1
	currentDirHoriz := 0
	stepDone := false

exit:
	for {
		if position.X == coords.X && position.Z == coords.Z {
			for y := position.Y - (extraChunksToLoadForEntity * int32(api.ChunkSideSize)); y <= position.Y+(extraChunksToLoadForEntity*int32(api.ChunkSideSize)); y += int32(api.ChunkSideSize) {
				coords,_ := utils.GetRightWorldCoords(api.WorldCoords{X: coords.X, Y: y, Z: coords.Z}, this.Universe.Size)

				c := this.getOrCreateChunk(coords)
				worldChunks[index] = c
				index++

				//log.Warnf(" ===== chunk %v / amountH %v / currentDirH %v / stepdone %v", c.Coords, amountHoriz, currentDirHoriz, stepDone)
			}
		}

		for i := 0; i < amountHoriz; i++ {
			switch currentDirHoriz {
			case 0:
				coords.X += int32(int32(api.ChunkSideSize))
			case 1:
				coords.Z += int32(int32(api.ChunkSideSize))
			case 2:
				coords.X -= int32(int32(api.ChunkSideSize))
			case 3:
				coords.Z -= int32(int32(api.ChunkSideSize))
			}

			if amountHoriz < chunksStorageLength || (!stepDone && i < amountHoriz-1) {
				for y := position.Y - (extraChunksToLoadForEntity * int32(api.ChunkSideSize)); y <= position.Y+(extraChunksToLoadForEntity*int32(api.ChunkSideSize)); y += int32(api.ChunkSideSize) {
					c := this.getOrCreateChunk(api.WorldCoords{X: coords.X, Y: y, Z: coords.Z})
					worldChunks[index] = c
					index++

					//log.Warnf(" ===== chunk %v / amountH %v / currentDirH %v / stepdone %v", c.Coords, amountHoriz, currentDirHoriz, stepDone)
				}
			}
		}

		if stepDone {
			amountHoriz++
			stepDone = false
		} else {
			if amountHoriz < chunksStorageLength {
				stepDone = true
			} else {
				amountHoriz++
			}
		}

		if amountHoriz > chunksStorageLength {
			break exit
			/*amountHoriz = 1
			amountVert++
			coords.X = position.X
			coords.Z = position.Z
			coords.Y += int32(amountVert * currentDirVert * int(int32(api.ChunkSideSize)))

			if currentDirVert == 1 {
				currentDirVert = -1
			} else {
				currentDirVert = 1
			}*/
		}

		currentDirHoriz++
		if currentDirHoriz > 3 {
			currentDirHoriz = 0
		}
	}

	for _, c := range worldChunks {
		if c != nil {
			chunkEvt := &events.Chunk{
				Coords: c.Coords,
				Chunk:  c,
				General: events.General{
					EventUniverseID: this.Universe.Name,
					EventTime:       this.Universe.Age,
					EventSource:     this.ID(),
					EventMetadata:   metadata,
				},
				Outputable: events.Outputable{
					EventTarget: source,
				},
			}

			this.deferredProcessResults <- chunkEvt
		}
	}

	log.Tracef("chunks %#v", len(worldChunks))

	return worldChunks
}

func (this *uniProcessorImpl) getOrCreateChunk(coord api.WorldCoords) *api.Chunk {
	log.Tracef("%s", coord)
	
	packed := coord.Pack()

	this.uniLock.Lock()
	
	this.Universe.Chunks.RLock()
	
	c, ok := this.Universe.Chunks.Data[packed]
	log.Debugf("asked %v => %v", coord, ok)
	this.Universe.Chunks.RUnlock();
	
	if !ok || c == nil {
		item := this.lru.Get(string(packed))
		if item != nil {
			c = item.Value().(*api.Chunk)
		}
	}

	if c == nil {
		this.Universe.Chunks.Lock()
		
		c = this.WorldGenerator.GenerateChunk(coord)
		this.lru.Set(string(packed), c, time.Minute * 30)
		this.Universe.Chunks.Data[packed] = c
		
		this.Universe.Chunks.Unlock();
	}

	this.uniLock.Unlock()

	return c
}

func (this *uniProcessorImpl) sendWorldChunks(we *events.World) {
	log.Tracef("%#v", we)

	this.entityChunks(we.World.Owner, we.Source(), int32(we.Radius), we.Metadata())
}

func (this *uniProcessorImpl) Process(time uint64, evts []api.Event) []api.Outputable {
	rends := make([]api.Outputable, 0, len(evts))

	for _, ev := range evts {
		if ev.Source() == this.ID() {
			continue
		}

		log.Debugf("%#v", ev)

		renderMe := true

		switch v := ev.(type) {
		case *events.World:
			world := this.World(v.World.Owner)
			v.World = world

			v.EventTarget = v.Source()
			v.EventSource = this.ID()

			if v.NeedChunks {
				go this.sendWorldChunks(v)
			}
		case *events.CustomBlocks:
			mapb := this.BlockRegistry.All()
			v.Blocks = make([]api.BlockDefinition, len(mapb))
			i := 0
			for _, bv := range mapb {
				v.Blocks[i] = bv
			}

			v.EventTarget = v.Source()
			v.EventSource = this.ID()

		case *events.Chunk:
			renderMe = false
			go func(v *events.Chunk) {
				c := this.getOrCreateChunk(v.Coords)
				v.Chunk = c

				v.EventTarget = v.Source()
				v.EventSource = this.ID()

				this.deferredProcessResults <- v
			}(v)
		case *events.Chunks:
			renderMe = false

			if v.Coords != nil && len(v.Coords) > 0 {
				go func(v *events.Chunks) {
					for i, crds := range v.Coords {
						c := this.getOrCreateChunk(crds)

						if v.SplitOutput {
							chunkEvt := &events.Chunk{
								Coords: crds,
								Chunk:  c,
								General: events.General{
									EventUniverseID: this.Universe.Name,
									EventTime:       this.Universe.Age,
									EventSource:     this.ID(),
									EventMetadata:   v.Metadata(),
								},
								Outputable: events.Outputable{
									EventTarget: v.Source(),
								},
							}

							this.deferredProcessResults <- chunkEvt
						} else {
							if v.Chunks == nil || len(v.Chunks) != len(v.Coords) {
								v.Chunks = make([]*api.Chunk, len(v.Coords))
							}

							v.Chunks[i] = c
						}
					}

					if !v.SplitOutput {
						vv := &events.Chunks{
							General:    v.General,
							Outputable: v.Outputable,
							Coords:     v.Coords,
							Chunks:     v.Chunks,
						}

						vv.EventTarget = v.Source()
						vv.EventSource = this.ID()

						this.deferredProcessResults <- vv
					}

				}(v)
			}
		case *events.EntityMove:
			renderMe = false

			if v.Offset.X != 0 || v.Offset.Y != 0 || v.Offset.Z != 0 {
				v.Offset = this.getEntityActualOffset(v.EntityID, v.Offset)

				if v.Offset.X != 0 || v.Offset.Y != 0 || v.Offset.Z != 0 || v.HorizontalAngle != 0 || v.VerticalAngle != 0 {
					pos := this.MoveEntity(v)
					if pos != nil {
						ev = pos
						v.EventSource = this.ID()
					}
				}
			}
		default:
			continue
		}

		if renderMe {
			if r, ok := ev.(api.Outputable); ok {
				//sev.SetSource(this.ID())
				rends = append(rends, r)
			} else {
				log.Warnf("not a renderable: %#v", ev)
			}
		}
	}

	select {
	case r := <-this.deferredProcessResults:
		log.Debugf("deferred %#v", r)
		rends = append(rends, r)
	default:
		//log.Tracef("uniProcessorImpl.Process: no defers")
	}

	return rends
}

func (this *uniProcessorImpl) Start() error {
	return nil
}

func (this *uniProcessorImpl) Stop() {
}

func (this *uniProcessorImpl) MoveEntity(event *events.EntityMove) *events.EntityPosition {
	entity := this.Entity(event.EntityID)
	if entity == nil {
		log.Warnf("no entity with ID# %s found", event.EntityID)
		return nil
	}

	player := this.players[event.EntityID]
	ep := player.MakeMove(event)

	//entity.Position = ep.Position
	//entity.LooksAt = ep.LooksAt

	return ep
}

func fillBlockEvents(deferredProcessResults chan api.Outputable, c *api.Chunk, name string, age uint64, source string, metadata interface{}) {
	log.Tracef("chunk %v, source %v", c.Coords.Pack(), source)

	index := 0
	var blocks *events.Blocks
	for y := 0; y < int(int32(api.ChunkSideSize)); y++ {
		for x := 0; x < int(int32(api.ChunkSideSize)); x++ {
			for z := 0; z < int(int32(api.ChunkSideSize)); z++ {
				if c.Blocks[x][y][z] != nil {
					if index == 0 {
						blocks = &events.Blocks{
							General: events.General{
								EventUniverseID: name,
								EventTime:       age,
								EventSource:     reflect.TypeOf((*modules.UniverseProcessor)(nil)).Elem().String(),
								EventMetadata:   metadata,
							},
							Outputable: events.Outputable{
								EventTarget: source,
							},
							Blocks: make(map[api.PackedWorldCoords]*api.Block, blockDataLengthForEvent),
						}
					}

					coords := api.WorldCoords{X: c.Coords.X + int32(x), Y: c.Coords.Y + int32(y), Z: c.Coords.Z + int32(z)}
					blocks.Blocks[coords.Pack()] = c.Blocks[x][y][z].Block
					index++
				}

				if blocks != nil && (index == blockDataLengthForEvent || (x == int(int32(api.ChunkSideSize)-1) && y == int(int32(api.ChunkSideSize)-1) && z == int(int32(api.ChunkSideSize)-1))) {
					//rends = append(rends, blocks)
					deferredProcessResults <- blocks

					index = 0
				}
			}
		}
	}
}
