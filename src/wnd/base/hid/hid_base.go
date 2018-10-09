package hid

import (
	"wnd/api"
)

type Key int

const (
	ActionID = "HIDAction"
	
	Forward Key = iota
	Back
	Left
	Right
	Up
	Down
	Exit
)

type Action struct {
	HorizontalAngle,VerticalAngle float64
	Offset api.Coords
	AcTime int64
}

func (this *Action) ID() string {
	return ActionID
}

func (this *Action) Time() uint64 {
	return uint64(this.AcTime)
}

func (this *Action) Source() string {
	return "HID"
}

func (this *Action) SetMetadata(source interface{}) {
}

func (this *Action) Metadata() interface{} {
	return nil
}