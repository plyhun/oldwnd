package worldprocess

import (
	"wnd/api"
	"wnd/api/events"
	"wnd/base/entity"
	"wnd/base/graphics"
	"wnd/base/hid"
	"wnd/modules"
	"wnd/utils"
	"wnd/utils/log"

	"reflect"
	"sync"

	systemtime "time"
)

const (
	_REQ_WORLD    = "world"
	_REQ_BLOCKS   = "blocks"
	_REQ_ENTITIES = "entities"
)

type worldProcessor struct {
	worldLock sync.Mutex

	requests               map[string]bool
	deferredProcessResults chan api.Outputable
	
	owner *entity.Playable

	World         *api.World            `inject:""`
	BlockRegistry modules.BlockRegistry `inject:""`
}

func New() modules.WorldProcessor {
	return &worldProcessor{
		requests:               make(map[string]bool),
		deferredProcessResults: make(chan api.Outputable),

		//lostBlocks: make([]*events.Blocks, 0),
	}
}

func (this *worldProcessor) ID() string {
	return "worldProcessor"
}

func (this *worldProcessor) Priority() int8 {
	return -10
}

func (this *worldProcessor) checkMissingData(time uint64) []api.Outputable {
	//log.Tracef("worldProcessor.checkMissingData")

	renderables := make([]api.Outputable, 0, 3)

	if !this.BlockRegistry.IsReady() {
		if _, ok := this.requests[_REQ_BLOCKS]; !ok {
			log.Tracef("requesting blocks")

			this.requests[_REQ_BLOCKS] = true
			renderables = append(renderables, this.requestBlocks(time))
		}
	}

	if this.World.Chunks.Data == nil || len(this.World.Chunks.Data) < 1 {
		if _, ok := this.requests[_REQ_WORLD]; !ok {
			log.Tracef("requesting chunks")

			this.requests[_REQ_WORLD] = true
			renderables = append(renderables, this.requestChunks(time))
		}
	}

	if this.World.Entities == nil || len(this.World.Entities) < 1 {
		if _, ok := this.requests[_REQ_ENTITIES]; !ok {
			log.Tracef("requesting entities")

			this.requests[_REQ_ENTITIES] = true
			//request entities
		}
	} else {
		delete(this.requests, _REQ_ENTITIES)
	}

	return renderables
}

func (this *worldProcessor) processChunk(c *api.Chunk) []*events.Chunk {
	log.Tracef("%v", c.Coords)

	out := make([]*events.Chunk, 0)
	
	if _, ok := this.World.Chunks.Data[c.Coords.Pack()]; !ok {
		this.worldLock.Lock()
		this.World.Chunks.Data[c.Coords.Pack()] = c
		this.worldLock.Unlock()

		if _, ok := this.requests[_REQ_WORLD]; ok {
			delete(this.requests, _REQ_WORLD)
		}
	} else {
		//c.BiomeData = v.BiomeData
	}

	return out
}

func (this *worldProcessor) Process(time uint64, evts []api.Event) []api.Outputable {
	renderables := make([]api.Outputable, 0, 10)

	for _, e := range evts {
		if e.Source() == this.ID() {
			continue
		}

		log.Tracef("%#v", e)

		switch v := e.(type) {
		case *events.World:
			this.World.Age = v.World.Age
			this.World.Name = v.World.Name
			this.World.Seed = v.World.Seed
			this.World.Size = v.World.Size
			this.World.Usage = graphics.TARGET
		case *events.CustomBlocks:
			this.BlockRegistry.RegisterBlocks(true, v.Blocks...)

			delete(this.requests, _REQ_BLOCKS)
		case *events.Chunks:
			if v.Chunks != nil {
				v.EventTarget = reflect.TypeOf((*modules.Graphics)(nil)).Elem().String()
				renderables = append(renderables, v)
				for _, c := range v.Chunks {
					//this.processChunk(c)

					for _, q := range this.processChunk(c) {
						renderables = append(renderables, q)
					}
				}
			}
		case *events.Chunk:

			//log.Warnf("%v", v.Chunk.Coords)

			v.EventTarget = reflect.TypeOf((*modules.Graphics)(nil)).Elem().String()
			renderables = append(renderables, v)

			for _, q := range this.processChunk(v.Chunk) {
				log.Warnf("update along %v", q.Coords)
				renderables = append(renderables, q)
			}
		case *events.Blocks:
			for k, b := range v.Blocks {
				coords, _ := api.UnpackWorldCoords(k)
				chunkCoords := coords.ChunkCoords()
				
				c, ok := this.World.Chunks.Data[chunkCoords.Pack()]

				if !ok {
					log.Warnf("no chunk found for %v", k)

					c = &api.Chunk{
						Coords: chunkCoords,
					}
					
					this.World.Chunks.Data[chunkCoords.Pack()] = c

					askChunk := &events.Chunk{
						General: events.General{
							EventSource:     this.ID(),
							EventID:         "event",
							EventTime:       v.Time(),
							EventUniverseID: this.World.ID(),
						},
						NoBlocks: true,
						Coords:   chunkCoords,
					}

					renderables = append(renderables, askChunk)
				} else {
					eChunk := &events.Chunk{
						General: events.General{
							EventSource:     this.ID(),
							EventTime:       v.Time(),
							EventUniverseID: this.World.ID(),
						},
						Coords: chunkCoords,
					}

					renderables = append(renderables, eChunk)
				}

				c.SetBlock(coords, b)
			}

			//renderables = append(renderables, this.World)
		case *events.EntityPosition:
			log.Debugf("position %v => %v", v.Position, v.LooksAt)

			//v.EventTarget = reflect.TypeOf((*modules.Graphics)(nil)).Elem().String()
			//renderables = append(renderables, v)
		case *hid.Action:
			if this.owner == nil || this.owner.Entity == nil {
				log.Errorf("!!! %v", this.owner)
				break
			}

			//log.Warnf("hid action %v (%v / %v) %v", v.Offset, v.HorizontalAngle, v.VerticalAngle, v.AcTime)

			r := &events.EntityMove{
				General: events.General{
					EventSource:     this.ID(),
					EventID:         "move",
					EventTime:       v.Time(),
					EventUniverseID: this.World.ID(),
				},
				Outputable: events.Outputable{
					EventTarget: reflect.TypeOf((*modules.UniverseProcessor)(nil)).Elem().String(),
				},
				DeltaTime:       v.AcTime,
				Offset:          v.Offset,
				HorizontalAngle: v.HorizontalAngle,
				VerticalAngle:   v.VerticalAngle,
				EntityID:        this.World.Owner,
			}

			oldPos := this.owner.Position

			ep := this.owner.MakeMove(r)

			//log.Warnf("move to %v / %v (%#v))", ep.Position, ep.LooksAt, r)

			var ok bool
			if ep.Position, ok = utils.GetRightCoords(ep.Position, this.World.Size); ok {
				ep.LooksAt, _ = utils.GetRightCoords(ep.LooksAt, this.World.Size)
			}

			this.owner.Position = ep.Position
			this.owner.LooksAt = ep.LooksAt

			ep.EventTarget = graphics.TARGET
			ep.IsPlayer = true

			renderables = append(renderables, r, ep)
			//renderables = append(renderables, this.World)

			oldChunkCoords, newChunkCoords := oldPos.ToWorldCoords().ChunkCoords(), ep.Position.ToWorldCoords().ChunkCoords()

			if oldChunkCoords.X != newChunkCoords.X || oldChunkCoords.Y != newChunkCoords.Y || oldChunkCoords.Z != newChunkCoords.Z {
				go this.processEntityChunks(oldChunkCoords, newChunkCoords, v.Time())
			}
		}
	}

	renderables = append(renderables, this.checkMissingData(time)...)

	//log.Debugf("renderables %v", renderables)
br:
	for {
		select {
		case r := <-this.deferredProcessResults:
			log.Debugf("deferred %#v", r)
			renderables = append(renderables, r)
		default:
			break br
			//log.Tracef("uniProcessorImpl.Process: no defers")
		}
	}

	return renderables
}

func (this *worldProcessor) getRequiredAndPurgeableChunks(oldChunkCoords, newChunkCoords api.WorldCoords) ([]api.WorldCoords, []api.WorldCoords) {
	oldChunks, newChunks := make([]api.WorldCoords, 0), make([]api.WorldCoords, 0)

	minx, miny, minz := newChunkCoords.X-int32(this.World.VisibleRadius*api.ChunkSideSize), newChunkCoords.Y-int32(this.World.VisibleRadius*api.ChunkSideSize), newChunkCoords.Z-int32(this.World.VisibleRadius*api.ChunkSideSize)
	maxx, maxy, maxz := newChunkCoords.X+int32((this.World.VisibleRadius+1)*api.ChunkSideSize), newChunkCoords.Y+int32((this.World.VisibleRadius+1)*api.ChunkSideSize), newChunkCoords.Z+int32((this.World.VisibleRadius+1)*api.ChunkSideSize)

	/*for x := minx; x <maxx; x += int32(api.ChunkSideSize) {
		for y := miny; y <maxy; y += int32(api.ChunkSideSize) {
			for z := minz; z <maxz; z += int32(api.ChunkSideSize) {
				coords := utils.GetRightWorldCoords(api.WorldCoords{X: int32(x), Y: int32(y), Z: int32(z)}, this.World.Size)
				oldChunks = append(oldChunks, coords)
			}
		}
	}*/
	
	this.World.Chunks.RLock()
	
	for _, c := range this.World.Chunks.Data {
		if c != nil {
			if c.Coords.X < minx || c.Coords.X >= maxx || c.Coords.Y < miny || c.Coords.Y >= maxy || c.Coords.Z < minz || c.Coords.Z >= maxz {
				crds, _ := utils.GetRightWorldCoords(c.Coords, this.World.Size)
				oldChunks = append(oldChunks, crds)
			}
		}
	}

	for x := minx; x < maxx; x += int32(api.ChunkSideSize) {
		for y := miny; y < maxy; y += int32(api.ChunkSideSize) {
			for z := minz; z < maxz; z += int32(api.ChunkSideSize) {
				//crds := api.WorldCoords{X: x, Y: y, Z: z}

				crds, _ := utils.GetRightWorldCoords(api.WorldCoords{X: int32(x), Y: int32(y), Z: int32(z)}, this.World.Size)

				if _, ok := this.World.Chunks.Data[crds.Pack()]; !ok {
					newChunks = append(newChunks, crds)
				}
			}
		}
	}
	
	this.World.Chunks.RUnlock()

	//log.Debugf("old %v, new %v", oldChunks, newChunks)

	return oldChunks, newChunks
}

func (this *worldProcessor) processEntityChunks(oldChunkCoords, newChunkCoords api.WorldCoords, time uint64) {
	//oldc, newc := utils.EntityGetRequiredAndPurgedChunks(oldChunkCoords, newChunkCoords, this.World.VisibleRadius)
	oldc, newc := this.getRequiredAndPurgeableChunks(oldChunkCoords, newChunkCoords)
	if len(newc) > 0 {
		for _, n := range newc {
			eChunk := &events.Chunk{
				General: events.General{
					EventSource:     this.ID(),
					EventTime:       time,
					EventUniverseID: this.World.ID(),
				},
				Coords: n,
			}

			this.World.Chunks.RLock()
			
			if c, ok := this.World.Chunks.Data[n.Pack()]; ok {
				eChunk.Chunk = c
				eChunk.Outputable = events.Outputable{
					EventTarget: graphics.TARGET,
				}
			} else {
				//log.Debugf("ask for %v", n)
				eChunk.Outputable = events.Outputable{
					EventTarget: reflect.TypeOf((*modules.UniverseProcessor)(nil)).Elem().String(),
				}
			}
			
			this.World.Chunks.RUnlock()
			this.deferredProcessResults <- eChunk
		}
	}

	if len(oldc) > 0 {
		for _, o := range oldc {
			go this.scheduleChunkDeletion(o)
		}
	}
}

func (this *worldProcessor) scheduleChunkDeletion(coords api.WorldCoords) {
	packed := coords.Pack()

	<-systemtime.Tick(systemtime.Second << 2)

	if !this.areEntitiesNearby(coords) {
		this.World.Chunks.Lock()
		defer this.World.Chunks.Unlock()
		
		if c, ok := this.World.Chunks.Data[packed]; ok {
			rem := &graphics.RemoveOutputable{
				OutputableID: graphics.GetGraphicsOutputableID(c),
			}

			this.deferredProcessResults <- rem
		}
		
		this.worldLock.Lock()
		delete(this.World.Chunks.Data, packed)
		this.worldLock.Unlock()

		log.Warnf("deleted %v", coords)
	} else {
		log.Warnf("delete %v failed, entities nearby")
	}
}

func (this *worldProcessor) areEntitiesNearby(coords api.WorldCoords) bool {
	if this.World.Entities == nil {
		return false
	}

	coords = coords.ChunkCoords()
	crds := coords.ToCoords()
	if shift,ok := this.World.GetShift(crds, this.owner.Position); ok {
		
		coords = api.Coords{X: crds.X+shift.X, Y: crds.Y+shift.Y, Z: crds.Z+shift.Z}.ToWorldCoords()
	}
	
	distance := int32(api.ChunkSideSize * this.World.VisibleRadius)

	for _, e := range this.World.Entities {
		entityChunk := e.Position.ToWorldCoords().ChunkCoords()

		if coords.Equals(entityChunk) {
			return true
		}

		dx := coords.X - entityChunk.X
		dy := coords.Y - entityChunk.Y
		dz := coords.Z - entityChunk.Z

		if dx < 0 {
			dx = -dx
		}
		if dy < 0 {
			dy = -dy
		}
		if dz < 0 {
			dz = -dz
		}

		if dx < distance && dy < distance && dz < distance {
			//log.Warnf(" ==== entity %v nearby %v", entityChunk, coords)
			return true
		}
	}

	return false
}

func (this *worldProcessor) requestBlocks(time uint64) api.Outputable {
	return &events.CustomBlocks{
		General: this.getGeneralEvent(time),
		Outputable: events.Outputable{
			EventTarget: reflect.TypeOf((*modules.UniverseProcessor)(nil)).Elem().String(),
		},
	}
}

func (this *worldProcessor) requestChunks(time uint64) api.Outputable {
	return &events.World{
		General: this.getGeneralEvent(time),
		Outputable: events.Outputable{
			EventTarget: reflect.TypeOf((*modules.UniverseProcessor)(nil)).Elem().String(),
		},
		World: &api.World{
			Owner: this.World.Owner,
		},
		Radius:     uint8(this.World.VisibleRadius),
		NeedChunks: true,
	}
}

func (this *worldProcessor) getGeneralEvent(time uint64) events.General {
	return events.General{
		EventSource:     this.ID(),
		EventID:         systemtime.Now().String(),
		EventTime:       time,
		EventUniverseID: this.World.ID(),
	}
}

func (this *worldProcessor) Start() error {
	this.owner = entity.NewPlayable(this.World.Entities[this.World.Owner])
	return nil
}

func (this *worldProcessor) Stop() {

}
