package main

import (
	"wnd/base"
	"wnd/base/runtime"
	"wnd/modules/graphicsgl"
	"wnd/modules/hidgl"
	"wnd/modules/networkreudp"
	"wnd/modules/worldprocess"
	"wnd/utils/log"
	"wnd/utils/test"
	
	//"github.com/davecheney/profile"
)

var (
	r base.Runtime
)

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()
	
	//log.Filter("FOV")
	log.LoggerLevel(log.LevelTrace)

	u := test.GetUniverse()
	ownerId := "player"
	r = runtime.NewUniverse(u, ownerId, false)

	r.AddModule(hidgl.NewMouseKb())
	r.AddModule(graphicsgl.New())

	r.AddModule(networkreudp.NewClient("localhost"))
	r.AddModule(worldprocess.New())

	r.Start()
}
