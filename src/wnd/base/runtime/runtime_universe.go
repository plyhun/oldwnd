package runtime

import (
	"wnd/api"
	"wnd/base"
	"wnd/modules/blockregistryfs"
	"wnd/modules/broadcastnetwork"
	"wnd/modules/gamedirfs"
	"wnd/modules/incomingnetwork"
	"wnd/modules/playerregistryfs"
	"wnd/modules/serializercodec"
	"wnd/modules/storagefs"
	"wnd/modules/worldgennoise"
	"wnd/utils/log"
	//"wnd/modules/autoconfiguratorcmd"
)

type UniverseRuntime struct {
	runtime
	*api.Universe
}

func NewUniverse(u *api.Universe, ownerId string, isServer bool) base.Runtime {
	log.Tracef("NewUniverse: universe %#v, owner %s, isServer=%v", u, ownerId, isServer)

	isClient := ownerId != ""

	r := &UniverseRuntime{
		runtime:  newRuntime(isServer, isClient),
		Universe: u,
	}

	r.injectArbitrary(u)
	r.injectArbitrary(r)

	if isClient {
		r.injectArbitrary(&api.World{Universe: u, Owner: ownerId, VisibleRadius: 4})
	}
	
	if isServer {
		r.AddModule(worldgennoise.New())
	}

	r.AddModule(gamedirfs.New())
	r.AddModule(serializercodec.New())
	r.AddModule(storagefs.New())
	r.AddModule(blockregistryfs.New())
	r.AddModule(playerregistryfs.New())
	//r.AddModule(autoconfiguratorcmd.New())
	
	if isClient != isServer {
		r.AddModule(broadcastnetwork.New())
		r.AddModule(incomingnetwork.New())
	}

	return r
}
