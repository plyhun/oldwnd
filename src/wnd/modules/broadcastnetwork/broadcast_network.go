package broadcastnetwork

import (
	"wnd/api"
	"wnd/api/events"
	"wnd/base"
	"wnd/modules"
	"wnd/utils/log"
	
	"strings"
	"reflect"
)

type broadcastNetwork struct {
	Network    modules.Network    `inject:""`
	Serializer modules.Serializer `inject:""`
	Runtime    base.Runtime       `inject:""`
}

func New() modules.Broadcast {
	return &broadcastNetwork{}
}

func (this *broadcastNetwork) ID() string {
	return "broadcastNetwork"
}

func (this *broadcastNetwork) SendTo(dst interface{}, e api.Event, reliable bool) error {
	log.Tracef("%v => %v", e.ID(), dst)

	if b, err := this.Serializer.Serialize(e); err != nil {
		return err
	} else {
		log.Debugf("sending as %v", len(b))
		return this.Network.Send(dst, b, reliable)
	}
}

func (this *broadcastNetwork) Output(time uint64, renderables []api.Outputable) {
	//log.Tracef("broadcastNetwork.Render: %v", renderables)

	go this.send(time, renderables)
}

func (this *broadcastNetwork) send(time uint64, renderables []api.Outputable) {
	for _, r := range renderables {
		//log.Tracef("%v of %v: %v (%v)", i, len(renderables), r.ID(), r.Target())
		
		if r.Target() == reflect.TypeOf((*modules.Graphics)(nil)).Elem().String() {
			log.Debugf("skipping graphics-targeted: %#v", r)
			continue
		}
		
		if evt, ok := r.(api.Event); ok {
			evtIdx := strings.Index(reflect.TypeOf(evt).String(), "events") 
			if evtIdx != 0 && evtIdx != 1 {
				log.Debugf("skipping non-broadcastable: %#v", evt)
				continue
			} 
			
			_, mayDrop := evt.(*events.EntityMove)

			to := evt.Metadata()

			switch {
			case !this.Runtime.IsServer():
				to = this.Network.Where()
			case to == nil:
				log.Debugf("%T (%+v) not suitable for broadcast")
				continue
			}

			err := this.SendTo(to, evt, !mayDrop)

			if err != nil {
				log.Errorf("%v", err)
			}
		}
	}
}

func (this *broadcastNetwork) Start() error {
	return nil
}

func (this *broadcastNetwork) Stop() {
}
