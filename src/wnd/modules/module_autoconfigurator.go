package modules

import (
	"wnd/api"
)

type AutoConfigurator interface {
	api.GameModule
	api.InittableModule
	api.PrioritizedModule
	
	Configure() error
}