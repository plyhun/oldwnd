package storagefs

import (
	"wnd/api"
	"wnd/modules"
	"wnd/utils/log"
	
	"fmt"
	"errors"
	"io/ioutil"
	"os"
	
	"github.com/snuk182/go-multierror"
)

type storageFs struct {
	GameDir modules.GameDir `inject:""`
	Serializer modules.Serializer `inject:""`
}

func New() modules.Storage {
	return &storageFs{}
}

func (this *storageFs) ID() string {
	return "storageFs"
}

func (this *storageFs) Priority() int8 {
	return 0
}

func (this *storageFs) NewUniverse(u *api.Universe) error {
	dir := this.GameDir.UniverseSavesDir(u.Name)
	
	var e *multierror.Error
	
	b,err := this.Serializer.Serialize(u)
	if (err != nil){
		e = multierror.Append(e, err)
	}
	
	e = multierror.Append(e, ioutil.WriteFile(dir + "/u.wndb", b, os.ModePerm))
	
	return e.ErrorOrNil()
}
func (this *storageFs) Universe(universeName string) *api.Universe {
	dir := this.GameDir.UniverseSavesDir(universeName)
	
	var e *multierror.Error
	
	b,err := ioutil.ReadFile(dir + "/u.wndb")
	if err != nil {
		e = multierror.Append(e, err)
	}
	
	i,err := this.Serializer.Deserialize(b)
	if u,ok := i.(*api.Universe); (!ok || err != nil) {
		if !ok {
			e = multierror.Append(e, errors.New(fmt.Sprintf("Cannot unparse %#v as *api.Universe", i)))
		}
		
		e = multierror.Append(e, err)
		
		log.Errorf("storageFs.Universe %v: error %v", universeName, e.ErrorOrNil())
		
		return nil
	} else {
		return u
	}
}
func (this *storageFs) DeleteUniverse(universeName string) {
	err := os.RemoveAll(this.GameDir.UniverseSavesDir(universeName))
	
	if err != nil {
		log.Errorf("storageFs.DeleteUniverse: %v", err)
	}
}
	
func (this *storageFs) SaveChunk(universeName string, chunk *api.Chunk) {
	dir := fmt.Sprintf("%s/terrain/%s", this.GameDir.UniverseSavesDir(universeName), chunk.Coords.Pack())
	
	var e *multierror.Error
	
	dto := &api.Chunk{
		BiomeData: api.BiomeData{
			Temperature: chunk.Temperature,
			Humidity: chunk.Humidity,
			HeightBias: chunk.HeightBias,
			SeaLevelHeight: chunk.SeaLevelHeight,
		},
	}
	
	b,err1 := this.Serializer.Serialize(dto)
	e = multierror.Append(e, err1)
	e = multierror.Append(e, ioutil.WriteFile(dir + "/c.wndb", b, os.ModePerm))
	
	blocks,err2 := this.Serializer.Serialize(chunk.Blocks)
	e = multierror.Append(e, err2)
	e = multierror.Append(e, ioutil.WriteFile(dir + "/b.wndb", blocks, os.ModePerm))
	
	err := e.ErrorOrNil()
	
	if err != nil {
		log.Errorf("storageFs.SaveChunk %v: %v", chunk.Coords.Pack(), err)
	}
}
func (this *storageFs) Chunk(universeName string, coords api.WorldCoords) (*api.Chunk,error) {
	dir := fmt.Sprintf("%s/terrain/%s", this.GameDir.UniverseSavesDir(universeName), coords.Pack())
	
	b1,err1 := ioutil.ReadFile(dir + "/c.wndb")
	
	if err1 != nil {
		return nil,err1
	}
	
	i1,err11 := this.Serializer.Deserialize(b1)
	if chunk,ok := i1.(*api.Chunk); (!ok || err11 != nil) {
		if !ok {
			return nil,errors.New(fmt.Sprintf("storageFs.Chunk %v: cannot cast %#v to *api.Chunk", coords.Pack(), i1))
		} else {
			return nil,err11
		}
	} else {
		var e *multierror.Error
	
		b2,err2 := ioutil.ReadFile(dir + "/b.wndb")
		
		e = multierror.Append(e, err2)
		
		i2,err22 := this.Serializer.Deserialize(b2)
		e = multierror.Append(e, err22)
		if blocks,ok := i2.([api.ChunkSideSize][api.ChunkSideSize][api.ChunkSideSize]*api.BlockInChunk); (!ok || err2 != nil) {
			if !ok {
				e = multierror.Append(e, errors.New(fmt.Sprintf("storageFs.Chunk %v: cannot cast %#v to [][][]*api.BlockData", coords.Pack(), i2)))
			} 
		} else {
			chunk.Blocks = blocks
		}
		
		return chunk, e.ErrorOrNil()
	}
}
	
func (this *storageFs) SaveEntity(universeName string, entity *api.Entity) {
	
}
func (this *storageFs) Entity(universeName, entityId string) (*api.Entity,error) {
	return nil,nil
}