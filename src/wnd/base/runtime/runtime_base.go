package runtime

import (
	"wnd/api"
	"wnd/utils/log"
	"wnd/modules/globalcontextmap"

	"reflect"
	"sort"
	"time"
	"unsafe"

	"github.com/facebookgo/inject"
	"github.com/snuk182/go-multierror"
)

type runtime struct {
	running bool
	isServer,isClient bool
	
	eventModules   []api.EventModule
	processModules []api.ProcessModule
	outputModules  []api.OutputModule
	otherModules   []api.GameModule

	pool inject.Graph
}

func newRuntime(isServer, isClient bool) runtime {
	log.Tracef("isServer %v / isClient %v", isServer, isClient)

	r := runtime{
		eventModules:   make([]api.EventModule, 0),
		processModules: make([]api.ProcessModule, 0),
		outputModules:  make([]api.OutputModule, 0),
		otherModules:   make([]api.GameModule, 0),
		running:        false,
		isServer: isServer,
		isClient: isClient,
	}
	
	r.AddModule(globalcontextmap.New())
	
	return r
}

func (this *runtime) IsClient() bool {
	return this.isClient
}
func (this *runtime) IsServer() bool {
	return this.isServer
}	

func (this *runtime) AddModule(m api.GameModule) error {
	var e *multierror.Error

	switch m.(type) {
	case api.EventModule:
		log.Tracef("%v as EventModule", m.ID())
		this.eventModules = append(this.eventModules, m.(api.EventModule))
		sort.Sort(sortmodules(*(*[]api.GameModule)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&this.eventModules))))))
	case api.ProcessModule:
		log.Tracef("%v as ProcessModule", m.ID())
		this.processModules = append(this.processModules, m.(api.ProcessModule))
		sort.Sort(sortmodules(*(*[]api.GameModule)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&this.processModules))))))
	case api.OutputModule:
		log.Tracef("%v as OutputModule", m.ID())
		this.outputModules = append(this.outputModules, m.(api.OutputModule))
		sort.Sort(sortmodules(*(*[]api.GameModule)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&this.outputModules))))))
	default:
		log.Tracef("%v", m.ID())
		this.otherModules = append(this.otherModules, m)
		sort.Sort(sortmodules(*(*[]api.GameModule)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&this.otherModules))))))
	}

	e = multierror.Append(e, this.pool.Provide(&inject.Object{Value: m}))

	if this.running {
		e = multierror.Append(e, this.pool.Populate())

		if im, ok := m.(api.InittableModule); ok {
			e = multierror.Append(e, im.Init())
		}
	}

	return e.ErrorOrNil()
}

func (this *runtime) injectArbitrary(a interface{}) error {
	log.Tracef("%#v", a)
	return this.pool.Provide(&inject.Object{Value: a})
}

func (this *runtime) initAll() error {
	log.Tracef("")

	var err *multierror.Error
	
	for _, m := range this.GetByType(reflect.TypeOf((*api.InittableModule)(nil)).Elem()) {
		im, _ := m.(api.InittableModule)
		err = multierror.Append(err, im.Init())
	}

	return err.ErrorOrNil()
}

func (this *runtime) RemoveModule(id string) {
	log.Tracef("%v", id)

	var m api.GameModule

	if mm := sortmodules(*(*[]api.GameModule)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&this.eventModules))))).remove(id); mm != nil {
		m = mm
	} else if mm := sortmodules(*(*[]api.GameModule)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&this.processModules))))).remove(id); mm != nil {
		m = mm
	} else if mm := sortmodules(*(*[]api.GameModule)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&this.outputModules))))).remove(id); mm != nil {
		m = mm
	} else if mm := sortmodules(*(*[]api.GameModule)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&this.otherModules))))).remove(id); mm != nil {
		m = mm
	}

	if m != nil {
		if dm, ok := m.(api.DestroyableModule); ok {
			dm.Destroy()
		}
	}
}

func (this *runtime) GetByType(t reflect.Type) []api.GameModule {
	log.Tracef("%v", t)

	results := make([]api.GameModule, 0)

	for _, v := range this.eventModules {
		if reflect.TypeOf(v).Implements(t) {
			results = append(results, v)
		}
	}

	for _, v := range this.processModules {
		if reflect.TypeOf(v).Implements(t) {
			results = append(results, v)
		}
	}

	for _, v := range this.outputModules {
		if reflect.TypeOf(v).Implements(t) {
			results = append(results, v)
		}
	}

	for _, v := range this.otherModules {
		if reflect.TypeOf(v).Implements(t) {
			results = append(results, v)
		}
	}
	
	return results
}

func (this *runtime) Start() error {
	log.Tracef("")

	var err *multierror.Error

	log.Debugf("%#v", this.pool)
	
	err = multierror.Append(err, this.pool.Populate())
	err = multierror.Append(err, this.initAll())
	
	//log.Infof("Inject: error %v, pool: %# v", err.ErrorOrNil(), this.pool.Objects())

	for _, m := range this.GetByType(reflect.TypeOf((*api.RuntimeModule)(nil)).Elem()) {
		rm, _ := m.(api.RuntimeModule)
		err = multierror.Append(err, rm.Start())
	}

	if err.ErrorOrNil() == nil {
		this.loop()
		
		//<-this.exit
	}

	return err.ErrorOrNil()
}

func (this *runtime) loop() {
	this.running = true

	for this.running {
		startLoop := uint64(time.Now().UnixNano())

		events := make([]api.Event, 0)
		renderables := make([]api.Outputable, 0)

		for _, em := range this.eventModules {
			e := em.Events(startLoop)
			events = append(events, e...)
		}

		for _, pm := range this.processModules {
			r := pm.Process(startLoop, events)
			renderables = append(renderables, r...)
		}

		for _, rm := range this.outputModules {
			rm.Output(startLoop, renderables)
		}

		//log.Infof("runtime: Game loop run: %d ns", uint64(time.Now().UnixNano())-startLoop)
	}
}

func (this *runtime) Stop() error {
	log.Tracef("")

	for _, m := range this.GetByType(reflect.TypeOf((*api.RuntimeModule)(nil)).Elem()) {
		rm, _ := m.(api.RuntimeModule)
		rm.Stop()
	}

	this.running = false
	this.destroy()
	
	return nil
}

func (this *runtime) destroy() {
	log.Tracef("")

	for _, m := range this.GetByType(reflect.TypeOf((*api.DestroyableModule)(nil)).Elem()) {
		rm, _ := m.(api.DestroyableModule)
		rm.Destroy()
	}
}

type sortmodules []api.GameModule

func (this sortmodules) remove(id string) api.GameModule {
	log.Tracef("%v", id)

	for i, v := range this {
		if v != nil && v.ID() == id {
			this = append(this[:i], this[i+1:]...)
			this = append([]api.GameModule(nil), this[:len(this)-1]...)
			return v
		}
	}

	return nil
}

func (this sortmodules) Len() int {
	return len(this)
}

func (this sortmodules) Less(i, j int) bool {
	ti, ok1 := this[i].(api.PrioritizedModule)
	tj, ok2 := this[j].(api.PrioritizedModule)

	if ok1 || ok2 {
		if ok1 && ok2 {
			return ti.Priority() < tj.Priority()
		} else {
			return ok1
		}
	} else {
		return false
	}
}
func (this sortmodules) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
