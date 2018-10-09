package eventlooper

import (
	"wnd/api"
	"wnd/modules"
	"wnd/utils/log"

	"math"
)

const (
	TYPE = "eventLooper"
)

func NewOut() api.OutputModule {
	return &eventLooperOut{}
}

type eventLooperOut struct {
	Context modules.GlobalContext `inject:""`
}

func (this *eventLooperOut) Priority() int8 {
	return math.MinInt8
}

func (this *eventLooperOut) ID() string {
	return "eventLooperOut"
}

func (this *eventLooperOut) Start() error {
	return nil
}

func (this *eventLooperOut) Stop() {
}

func (this *eventLooperOut) Output(time uint64, renderables []api.Outputable) {
	data := make([]api.Event, 0)
	o := this.Context.Get(TYPE)
	
	var obj []api.Event
	if o == nil {
		obj = make([]api.Event, 0)
	} else {
		obj = o.([]api.Event)
	}

	for _, v := range renderables {
		if vv, ok := v.(api.Event); ok {
			stale := false
			
			for _,c := range obj {
				if c != nil && c == vv {
					stale = true
					break
				}
			}
			
			if !stale {
				log.Debugf("%#v", vv)
				data = append(data, vv)
			} else {
				//log.Warnf("stale %#v", vv)
			}
		}
	}

	if len(data) > 0 {
		this.Context.Put(TYPE, data)
	} else {
		this.Context.Delete(TYPE)
	}
}

func NewIn() api.EventModule {
	return &eventLooperIn{}
}

type eventLooperIn struct {
	Context modules.GlobalContext `inject:""`
}

func (this *eventLooperIn) Priority() int8 {
	return math.MaxInt8
}

func (this *eventLooperIn) ID() string {
	return "eventLooperIn"
}

func (this *eventLooperIn) Start() error {
	return nil
}

func (this *eventLooperIn) Stop() {
}

func (this *eventLooperIn) Events(time uint64) []api.Event {
	obj := this.Context.Get(TYPE)
	
	if obj != nil {
		//defer this.Context.Delete(TYPE)
		
		if data, ok := obj.([]api.Event); ok {
			if len(data) > 0 {
				log.Debugf("%#v", data)
			}
			return data
		}
	}

	return make([]api.Event, 0)
}
