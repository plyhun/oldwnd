package main

import (
	"wnd/base"
	"wnd/base/runtime"
	"wnd/modules/eventlooper"
	"wnd/modules/graphicsgl"
	"wnd/modules/hidgl"
	"wnd/modules/uniprocess"
	"wnd/modules/worldprocess"
	"wnd/utils/log"
	"wnd/utils/test"
	
	//"github.com/davecheney/profile"
)

var (
	r base.Runtime
)

func main() {
	/*cfg := &profile.Config{
        MemProfile:     true,
        CPUProfile: true,
        BlockProfile: true,
        ProfilePath:    ".",  // store profiles in current directory
    }
	defer profile.Start(cfg).Stop()*/
	
	log.LoggerLevel(log.LevelInfo)

	ownerId := "player"
	u := test.GetUniverse()
	r = runtime.NewUniverse(u, ownerId, true)

	r.AddModule(worldprocess.New())
	r.AddModule(uniprocess.New())

	r.AddModule(hidgl.NewMouseKb())
	r.AddModule(graphicsgl.New())

	r.AddModule(eventlooper.NewIn())
	r.AddModule(eventlooper.NewOut())

	r.Start()
}
