package blockregistryfs

import (
	"wnd/api"
	"wnd/modules"
	"wnd/utils"
	"wnd/utils/log"

	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
)

const (
	_FORMAT_BLOCK = "%s/%s/block.json"
	_BLOCKS_DIR = "blocks"
)

type blocksRegistryFs struct {
	blocks map[string]api.BlockDefinition

	GameDir modules.GameDir `inject:""`
	Universe *api.Universe `inject:""`
}

func New() modules.BlockRegistry {
	return new(blocksRegistryFs)
}

func (this *blocksRegistryFs) All() map[string]api.BlockDefinition {
	log.Tracef("blocksRegistryFs.All: %#v", this)
	return this.blocks
}

func (this *blocksRegistryFs) ID() string {
	return "blockRegistryFs"
}

func (this *blocksRegistryFs) IsReady() bool {
	return this.blocks != nil && len(this.blocks) > 0
}

func (this *blocksRegistryFs) Priority() int8 {
	return -116
}

func (this *blocksRegistryFs) Init() error {
	log.Tracef("blocksRegistryFs.Init")

	this.blocks = make(map[string]api.BlockDefinition)

	dir1 := fmt.Sprintf("%s/resources/%s/", this.GameDir.AppDir(), _BLOCKS_DIR)
	
	return this.initFromDir(dir1)
}

func (this *blocksRegistryFs) initFromDir(dir string) error {
	log.Tracef("%v", dir)

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Errorf("Cannot read local block storage: %v", err)
		return err
	}

	for _, d := range files {
		if d.IsDir() {
			bytes, err := ioutil.ReadFile(fmt.Sprintf(_FORMAT_BLOCK, dir, d.Name()))
			if err != nil {
				log.Warnf("Cannot read block data %s: %v", d.Name(), err)
				continue
			}

			b, err := this.parseBlockJson(bytes, d.Name())
			if err != nil {
				log.Warnf("Cannot read block data %s: %v", d.Name(), err)
				continue
			}

			this.RegisterBlocks(false, b)
		} else {
			log.Warnf("Alien file in local block storage: %v", d.Name())
		}
	}

	return nil
}

func (this *blocksRegistryFs) RegisterBlocks(persist bool, blocks ...api.BlockDefinition) {
	log.Tracef("%#v", blocks)

	for _, b := range blocks {
		_, ok := this.blocks[b.Type]
		if b.ID != 0 && ok {
			log.Warnf("Cannot register %v: same typed block already exists", b.ID)
		} else {
			b.ID = uint32(len(this.blocks))
			this.blocks[b.Type] = b
			log.Infof("Block %v registered as ID# %v", b.Type, b.ID)
			this.storeBlock(b)
		}
	}
}

func (this *blocksRegistryFs) storeBlock(b api.BlockDefinition) {
	dir := this.getUniverseSavesDir() + "/" + b.Type
	
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		log.Warnf("Cannot store %s - block already exists", b.Type)
		return
	} 
	
	if err := utils.CheckAndMakeDir(dir); err == nil {
		if bytes,err1 := json.Marshal(b); err1 == nil {
			ioutil.WriteFile(fmt.Sprintf(_FORMAT_BLOCK, this.getUniverseSavesDir(), b.Type), bytes, 0777)
		} else {
			log.Errorf("%v",err1)
		}
	} else {
		log.Errorf("%v",err)
	}
}

func (this *blocksRegistryFs) UnregisterBlocks(blocks ...string) {
	for _, b := range blocks {
		_, ok := this.blocks[b]
		if ok {
			delete(this.blocks, b)
			log.Infof("Block %v unregistered", b)
		}
	}
}

func (this *blocksRegistryFs) getUniverseSavesDir() string {
	return this.GameDir.UniverseSavesDir(this.Universe.ID()) + "/" + _BLOCKS_DIR
}

func (this *blocksRegistryFs) AddUniverse() error {
	dir1 := this.GameDir.UniverseSavesDir(this.Universe.ID())
	return this.initFromDir(dir1)
}

func (this *blocksRegistryFs) DeleteUniverse() error {
	dir1 := this.GameDir.UniverseSavesDir(this.Universe.ID())
	return os.RemoveAll(dir1)
}	

func (this *blocksRegistryFs) Dir(b api.BlockDefinition) string {
	return fmt.Sprintf("%s/resources/%s/%s", this.GameDir.AppDir(), _BLOCKS_DIR, b.Type)
}

func (this *blocksRegistryFs) ByID(id uint32) api.BlockDefinition {
	for _,b := range this.blocks {
		if b.ID == id {
			return b
		}
	}
	
	return api.BlockDefinition{}
}

func (this *blocksRegistryFs) ByType(blockType string) api.BlockDefinition {
	return this.blocks[blockType]
}

func (this *blocksRegistryFs) Statistics() (int, string) {
	s := ""

	for _, b := range this.blocks {
		s += string(b.ID)
		s += " "
	}

	return len(this.blocks), s
}

func (this *blocksRegistryFs) parseBlockJson(jsonFile []byte, btype string) (b api.BlockDefinition, err error) {
	err = json.Unmarshal(jsonFile, &b)
	b.Type = btype
	return
}
