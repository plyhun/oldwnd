package modules

import (
	"wnd/api"
)

type PlayerRegistry interface {
	api.GameModule
	
	PlayerByID(playerId string) *api.Player
	PlayerByEntityID(entityID string) *api.Player
	PlayerByIRL(IRL interface{}) *api.Player
	AddPlayer(player *api.Player) error
	DeletePlayer(playerId string)
}