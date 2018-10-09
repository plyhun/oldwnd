package modules

import (
	"wnd/api"
)

type BlockRegistry interface {
	api.GameModule
	api.InittableModule
	api.PrioritizedModule
	
	Dir(b api.BlockDefinition) string
	ByID(id uint32) api.BlockDefinition
	ByType(blockType string) api.BlockDefinition
	RegisterBlocks(persist bool, blocks ...api.BlockDefinition)
	All() map[string]api.BlockDefinition
	Statistics() (int, string)
	IsReady() bool
}
