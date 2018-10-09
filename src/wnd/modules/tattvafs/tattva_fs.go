package tattvafs

import (
	"wnd/api"
	"wnd/modules"
	
	"errors"
	"fmt"
)

type tattvaFs struct {
	list []*api.Universe
	
	Storage modules.Storage `inject:""`
}

func New() modules.Tattva {
	return &tattvaFs{list: make([]*api.Universe, 0)}
}

func (this *tattvaFs) ID() string {
	return "tattva"
}

func (this *tattvaFs) NewUniverse(u *api.Universe) error {
	if this.HasUniverse(u.Name, true) {
		return errors.New(fmt.Sprintf("tattvaFs.NewUniverse: universe named '%s' already exists", u.Name))
	}
	
	this.list = append(this.list, u)
	return this.Storage.NewUniverse(u)
}

func (this *tattvaFs) UniverseList() []*api.Universe {
	return this.list
}

func (this *tattvaFs) Universe(universeId string) *api.Universe {
	for _,u := range this.list {
		if u != nil && u.Name == universeId {
			return u
		}
	}
	
	u := this.Storage.Universe(universeId)
	if u != nil {
		this.list = append(this.list, u)
	}
	
	return u
}

func (this *tattvaFs) HasUniverse(universeId string, stored bool) bool {
	if stored {
		return this.Storage.Universe(universeId) != nil
	} else {
		for _,u := range this.list {
			if u != nil && u.Name == universeId {
				return true
			}
		}
	}
	
	return false
}

func (this *tattvaFs) RemoveUniverse(universeId string, fromStorage bool) {
	for i,u := range this.list {
		if u != nil && u.Name == universeId {
			this.list = append(this.list[:i], this.list[i+1:]...)
			this.list = append([]*api.Universe(nil), this.list[:len(this.list)-1]...)
			break
		}
	}
	
	if fromStorage {
		this.Storage.DeleteUniverse(universeId)
	}
}