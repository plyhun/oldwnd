package modules

import (
	"wnd/api"
)

type GameDir interface {
	api.GameModule
	api.PrioritizedModule
	
	AppDir() string
	StorageDir() string
	SavesDir() string
	UniverseSavesDir(worldName string) string
}