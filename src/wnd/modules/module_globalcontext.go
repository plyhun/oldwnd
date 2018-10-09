package modules

import (
	"wnd/api"
)

type GlobalContext interface {
	api.GameModule
	
	Put(key string, value interface{}) error
	Get(key string) interface{}
	Delete(key string) interface{}
}