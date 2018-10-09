package globalcontextmap

import (
	"wnd/modules"
	"wnd/utils/log"
	
	"sync"
)

type globalContextMap struct {
	sync.RWMutex
	
	context map[string]interface{}
}

func New() modules.GlobalContext {
	return &globalContextMap{
		context: make(map[string]interface{}),
	}
}

func (this *globalContextMap) ID() string {
	return "globalContext"
}
func (this *globalContextMap) Put(key string, value interface{}) error {
	log.Tracef("%s = %v", key, value)
	
	this.Lock()
	defer this.Unlock()
	
	this.context[key] = value
	return nil
}
func (this *globalContextMap) Get(key string) interface{} {
	log.Tracef("%s from %v", key, this.context)
	
	this.Lock()
	defer this.Unlock()
	
	return this.context[key]
}
func (this *globalContextMap) Delete(key string) interface{} {
	log.Tracef("globalContextMap.Delete: %s from %v", key, this.context)
	
	this.Lock()
	defer this.Unlock()
	
	value,ok := this.context[key]
	delete(this.context, key)
	
	if ok {
		return value
	} else {
		return nil
	}
}