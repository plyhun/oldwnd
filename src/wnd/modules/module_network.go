package modules

import (
	"wnd/api"
)

type Network interface {
	api.RuntimeModule
	api.ConfigurableModule
	
	Send(where interface{}, b []byte, reliable bool) error
	Poll() ([]interface{}, [][]byte, error)
	Where() interface{}
}