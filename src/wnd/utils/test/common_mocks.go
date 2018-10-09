package test

import (
	"wnd/api"
	"wnd/modules"
	"wnd/modules/blockregistryfs"
	"wnd/utils"
	"wnd/utils/log"

	"math"
	"testing"

	"github.com/facebookgo/inject"
)

type GameModuleMock struct{}

func (this *GameModuleMock) Priority() int8 {
	return 0
}
func (this *GameModuleMock) ID() string {
	return "Mock"
}

type InittableModuleMock struct{}

func (this *InittableModuleMock) Init() error {
	return nil
}

//RuntimeModule 

type RuntimeModuleMock struct{
	GameModuleMock
}

func (this *RuntimeModuleMock) Start() error {
	return nil
}
func (this *RuntimeModuleMock) Stop() {
	
}

//GameDir

type gameDirMock struct {
	*InittableModuleMock
}

func (this *gameDirMock) AppDir() string {
	return "c:/opengl-libs"
}

func (this *gameDirMock) StorageDir() string {
	return "c:/opengl-libs"
}

func (this *gameDirMock) SavesDir() string {
	dir := this.StorageDir() + "/saves"
	
	err := utils.CheckAndMakeDir(dir)
	if err != nil {
		log.Errorf("gamedirFs.StorageDir: %v", err)
	}
	
	return dir
}

func (this *gameDirMock) UniverseSavesDir(worldName string) string {
	dir := this.SavesDir() + "/" + worldName
	
	err := utils.CheckAndMakeDir(dir)
	if err != nil {
		log.Errorf("gamedirFs.StorageDir: %v", err)
	}
	
	return dir
}

func (this *gameDirMock) Priority() int8 {
	return math.MinInt8
}

func (this *gameDirMock) ID() string {
	return "GameDirMock"
}

func NewGameDir() modules.GameDir {
	return &gameDirMock{}
}

// BlockRegistry

func NewBlockRegistry(t *testing.T) modules.BlockRegistry {
	var pool inject.Graph

	br := blockregistryfs.New()
	pool.Provide(&inject.Object{Value: br}, &inject.Object{Value: NewGameDir()})

	t.Logf("NewBlockRegistry: err %v", pool.Populate())

	return br
}

func GetUniverse() *api.Universe {
	return &api.Universe{
		Name: "genesis", 
		Seed: 8, 
		Size: 16, 
		Age: 12345678, 
		Chunks: api.Chunks{Data: make(map[api.PackedWorldCoords]*api.Chunk)}, 
		Entities: map[string]*api.Entity{
			"player": &api.Entity{
				ID: "player", 
				Control: "DaPlayer", 
				Speed: 25,
				Observer: api.Observer{
					Position: api.Coords{X: 20, Y: 30, Z: 40}, 
					LooksAt: api.Coords{X: 10, Y: 11, Z: 12},
				},
			},
		},
	}
}
