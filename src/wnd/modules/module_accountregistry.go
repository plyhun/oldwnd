package modules

import (
	"wnd/api"
)

type AccountRegistry interface {
	api.GameModule
	
	Login(id, password string, createIfAbsent bool) (*api.Account,error)
	Logout(id string) error
	Delete(id string) error
	Get(id string) (*api.Account,error)
	Update(a *api.Account) error
	UpdatePassword(id, oldpw, newpw string) error
}