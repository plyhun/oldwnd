package main

import (
	"wnd/base"
	"wnd/base/runtime"
	"wnd/modules/networkreudp"
	"wnd/modules/uniprocess"
	"wnd/utils/log"
	"wnd/utils/test"
)

var (
	runtimes []base.Runtime
)

func main() {
	log.LoggerLevel(log.LevelInfo)

	u := runtime.NewUniverse(test.GetUniverse(), "", true)

	u.AddModule(uniprocess.New())
	u.AddModule(networkreudp.NewServer())

	addRuntime(u)
}

func addRuntime(r base.Runtime) {
	runtimes = append(runtimes, r)

	r.Start()
}

func removeRuntime(id string) base.Runtime {
	for i, v := range runtimes {
		if v != nil && v.ID() == id {
			runtimes = append(runtimes[:i], runtimes[i+1:]...)
			runtimes = append([]base.Runtime(nil), runtimes[:len(runtimes)-1]...)

			if err := v.Stop(); err != nil {
				log.Errorf("server.removeRuntime: %v", err)
			}

			return v
		}
	}

	return nil
}
