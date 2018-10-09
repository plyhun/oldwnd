package runtime

import (
	"wnd/base"
	"wnd/modules"
	"wnd/modules/tattvafs"
	
	"sync"
)

var (
	instance *TattvaRuntime
	mutex *sync.Mutex
)

type TattvaRuntime struct {
	runtime
	modules.Tattva
}

func (this *TattvaRuntime) ID() string {
	return ""
}

func GetTattwa() base.Runtime {
	if instance == nil {
		mutex.Lock()
		if instance == nil {
			instance = &TattvaRuntime{
				runtime: newRuntime(false, false),
				Tattva: tattvafs.New(),
			}
			
			instance.AddModule(instance.Tattva)
			instance.injectArbitrary(instance)
		}
	}
	
	return instance
}