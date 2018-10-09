package gamedirfs

import (
	"wnd/api"
	"wnd/modules"
	"wnd/utils"
	"wnd/utils/log"

	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const (
	_PARAM_GAMEDIR_STORAGE = "storage"
	_PARAM_GAMEDIR_SAVES   = "saves"
)

var (
	storageDir = api.TypeKeyValue{
		Name:  "Storage directory",
		Key:   _PARAM_GAMEDIR_STORAGE,
		Type:  reflect.String,
		Value: "",
	}
	savesDir = api.TypeKeyValue{
		Name:  "Saves directory",
		Key:   _PARAM_GAMEDIR_SAVES,
		Type:  reflect.String,
		Value: "",
	}
)

func New() modules.GameDir {
	log.Tracef("gamedirFs.New")
	return &gamedirFs{storageDir: "", savesDir: ""}
}

type gamedirFs struct {
	storageDir string
	savesDir   string
}

func (this *gamedirFs) Priority() int8 {
	return -128
}

func (this *gamedirFs) ID() string {
	return "gamedirFs"
}

func (this *gamedirFs) AppDir() string {
	return utils.GetAppDir()
}

func (this *gamedirFs) StorageDir() string {
	var dir string
	var err error

	if this.storageDir == "" {
		dir, err = homedir.Expand("~/.wnd")
		if err != nil {
			log.Errorf("gamedirFs.StorageDir: %v", err)
		}
	} else {
		dir = this.storageDir
	}

	log.Tracef("gamedirFs.StorageDir: %#v", dir)

	dir = strings.Replace(dir, "\\", "/", -1)

	err = utils.CheckAndMakeDir(dir)
	if err != nil {
		log.Errorf("gamedirFs.StorageDir: %v", err)
	}

	return dir
}

func (this *gamedirFs) SavesDir() string {
	var dir string

	if this.savesDir == "" {
		dir = this.StorageDir() + "/saves"
	} else {
		dir = this.savesDir
	}

	log.Tracef("gamedirFs.SavesDir: %#v", dir)

	dir = strings.Replace(dir, "\\", "/", -1)

	err := utils.CheckAndMakeDir(dir)
	if err != nil {
		log.Errorf("gamedirFs.StorageDir: %v", err)
	}

	return dir
}

func (this *gamedirFs) UniverseSavesDir(worldName string) (dir string) {
	dir = fmt.Sprintf("%s/%s", this.SavesDir(), worldName)

	log.Tracef("gamedirFs.UniverseSavesDir: %#v", dir)

	err := utils.CheckAndMakeDir(dir)
	if err != nil {
		log.Errorf("gamedirFs.StorageDir: %v", err)
	}

	return
}

func (this *gamedirFs) SetConfiguration(values ...api.TypeKeyValue) error {
	log.Tracef("gamedirFs: configuration changed: %#v", values)

	for _, v := range values {
		switch v.Key {
		case _PARAM_GAMEDIR_STORAGE:
			this.storageDir = v.Value.(string)
		case _PARAM_GAMEDIR_SAVES:
			this.savesDir = v.Value.(string)
		}
	}

	return nil
}

func (this *gamedirFs) Configuration() (c []api.TypeKeyValue) {
	savesDir.Value = this.savesDir
	storageDir.Value = this.StorageDir()
	c = []api.TypeKeyValue{storageDir, savesDir}

	log.Tracef("gamedirFs.Configuration: %#v", c)

	return
}
