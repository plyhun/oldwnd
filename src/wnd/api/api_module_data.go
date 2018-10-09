package api

import (
	"reflect"
)

type Positionnable interface {
	Position() Coords
	SetPosition(position Coords)
}

type Modifiable interface {
	isModified() bool
}

type Event interface {
	ID() string
	Time() uint64 
	Source() string
	Metadata() interface{}
	SetMetadata(src interface{})
}

type Outputable interface {
	ID() string
	Target() string
}

type TypeKeyValue struct {
	Name LocalizableString
	Type reflect.Kind
	Key string
	Value interface{}
}

func (this TypeKeyValue) Clone() TypeKeyValue {
	return this
}