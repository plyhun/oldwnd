package modules

import (
	"wnd/api"
)

type Graphics interface {
	api.OutputModule
	
	AddToRender(r api.Outputable)
	//FOV(camera graphics.Camera)
	RemoveFromRender(rid string)
}