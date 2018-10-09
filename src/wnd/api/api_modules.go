package api

type GameModule interface {
	ID() string
}

type RuntimeModule interface {
	GameModule

	Start() error
	Stop()
}

type PrioritizedModule interface {
	Priority() int8
}

type ConfigurableModule interface {
	Configuration() []TypeKeyValue
	SetConfiguration(values ...TypeKeyValue) error
}

type InittableModule interface {
	Init() error
}

type DestroyableModule interface {
	Destroy()
}

type EventModule interface {
	RuntimeModule
	Events(time uint64) []Event
}

type ProcessModule interface {
	RuntimeModule
	PrioritizedModule
	Process(time uint64, events []Event) []Outputable
}

type OutputModule interface {
	RuntimeModule
	Output(time uint64, renderables []Outputable)
}

/*type InitModule interface {
	GameModule
	InittableModule
}

type DestroyModule interface {
	GameModule
	DestroyableModule
}*/
