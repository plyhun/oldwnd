package playerregistryfs

import (
	"wnd/api"
	"wnd/modules"
	"wnd/utils/log"
	
	"io/ioutil"
	"errors"
	"fmt"
	"os"
	
	"github.com/snuk182/go-multierror"
)

type playerRegistryFs struct {
	GameDir modules.GameDir `inject:""`
	Serializer modules.Serializer `inject:""`
	
	players map[string]*api.Player
}

func New() modules.PlayerRegistry {
	return &playerRegistryFs{
		players: make(map[string]*api.Player),
	}
}

func (this *playerRegistryFs) ID() string {
	return "playerRegistryFs"
}

func (this *playerRegistryFs) PlayerByID(playerId string) *api.Player {
	if p,ok := this.players[playerId]; ok {
		return p
	} else {
		return this.readFromStorage(playerId)
	}
}

func (this *playerRegistryFs) PlayerByEntityID(entityID string) *api.Player {
	for _,p := range this.players {
		if p.EntityID == entityID {
			return p
		}
	}
	
	return nil
}

func (this *playerRegistryFs) PlayerByIRL(IRL interface{}) *api.Player {
	for _,p := range this.players {
		if p.IRL == IRL {
			return p
		}
	}
	
	return nil
}

func (this *playerRegistryFs) AddPlayer(player *api.Player) error {
	file := this.GameDir.StorageDir() + "/players/" + player.EntityID + ".wndb"
	
	var e *multierror.Error
	
	b,err := this.Serializer.Serialize(player)
	if (err != nil){
		e = multierror.Append(e, err)
	}
	
	e = multierror.Append(e, ioutil.WriteFile(file, b, os.ModePerm))
	
	return e.ErrorOrNil()
}

func (this *playerRegistryFs) DeletePlayer(playerId string) {
	err := os.RemoveAll(this.GameDir.StorageDir() + "/players/" + playerId + ".wndb")
	
	if err != nil {
		log.Errorf("storageFs.DeleteUniverse: %v", err)
	}
}

func (this *playerRegistryFs) readFromStorage(playerId string) *api.Player {
	file := this.GameDir.StorageDir() + "/players/" + playerId + ".wndb"
	
	var e *multierror.Error
	
	b,err := ioutil.ReadFile(file)
	if err != nil {
		e = multierror.Append(e, err)
	}
	
	i,err := this.Serializer.Deserialize(b)
	if p,ok := i.(*api.Player); (!ok || err != nil) {
		if !ok {
			e = multierror.Append(e, errors.New(fmt.Sprintf("Cannot unparse %#v as *api.Player", i)))
		}
		
		e = multierror.Append(e, err)
		
		log.Errorf("playerRegistryFs.PlayerByID %v: error %v", playerId, e.ErrorOrNil())
		
		return nil
	} else {
		return p
	}
}