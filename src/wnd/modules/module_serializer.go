package modules

import (
	"wnd/api"
)

type Serializer interface {
	api.GameModule
	api.InittableModule
	api.DestroyableModule
	api.PrioritizedModule
	
	Serialize(i interface{}) ([]byte,error)
	Deserialize(b []byte) (interface{}, error)
}