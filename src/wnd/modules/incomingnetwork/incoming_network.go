package incomingnetwork

import (
	"wnd/api"
	"wnd/modules"
	"wnd/utils/log"
)

type incomingNet struct {
	Network    modules.Network    `inject:""`
	Serializer modules.Serializer `inject:""`
}

func New() modules.Incoming {
	return &incomingNet{}
}

func (this *incomingNet) ID() string {
	return "incomingNetwork"
}

func (this *incomingNet) Start() error {
	return nil
}

func (this *incomingNet) Stop() {
}

func (this *incomingNet) Parse(from,what interface{}) api.Event {
	log.Tracef("from %v -> %v", from,what)
	
	if bytes,ok := what.([]byte); ok {
		return this.parse(from, bytes)
	} else {
		log.Warnf("cannot parse %v = not a []byte", what)
		return nil
	}
}

func (this *incomingNet) parse(from interface{}, what []byte) api.Event {
	if w,e1 := this.Serializer.Deserialize(what); e1 == nil {
		if evt,ok := w.(api.Event); ok {
			evt.SetMetadata(from)
			return evt
		} else {
			log.Errorf("incomingNet.Events: not an event: %#v", what)
		}
	} else {
		log.Errorf("incomingNet.Events: %v", e1)
	}
	
	return nil
}

func (this *incomingNet) Events(time uint64) []api.Event {
	if froms, whats, e := this.Network.Poll(); e == nil {
		if len(froms) == len(whats) {
			events := make([]api.Event, 0, len(froms))
			
			for i := 0; i < len(froms); i++ {
				if evt := this.parse(froms[i], whats[i]); evt != nil {
					events = append(events, evt)
				}
			}
			
			return events
		} else {
			log.Errorf("Different results: %d/%d", len(froms), len(whats))
		}
	} else {
		log.Errorf("%v", e)
	}
	
	return []api.Event{}
}