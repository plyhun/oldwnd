package modules

import "wnd/api"

type HID interface {
	api.EventModule
	api.ConfigurableModule
}
