package hidgl

import (
	"wnd/api"
	"wnd/base"
	"wnd/base/hid"
	"wnd/modules"
	"wnd/utils/log"

	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-gl/glfw3/v3.1/glfw"
	"github.com/snuk182/go-multierror"
)

const (
	_KEY_FORWARD = "keyForward"
	_KEY_BACK    = "keyBack"
	_KEY_LEFT    = "keyLeft"
	_KEY_RIGHT   = "keyRight"
	_KEY_UP      = "keyUp"
	_KEY_DOWN    = "keyDown"
	_KEY_EXIT    = "keyExit"
)

var (
	keyForward = api.TypeKeyValue{Type: reflect.Int, Key: _KEY_FORWARD, Name: api.LocalizableString("'Forward' move key"), Value: glfw.KeyW}
	keyBack    = api.TypeKeyValue{Type: reflect.Int, Key: _KEY_BACK, Name: api.LocalizableString("'Back' move key"), Value: glfw.KeyS}
	keyLeft    = api.TypeKeyValue{Type: reflect.Int, Key: _KEY_LEFT, Name: api.LocalizableString("'Left' move key"), Value: glfw.KeyA}
	keyRight   = api.TypeKeyValue{Type: reflect.Int, Key: _KEY_RIGHT, Name: api.LocalizableString("'Right' move key"), Value: glfw.KeyD}
	keyUp      = api.TypeKeyValue{Type: reflect.Int, Key: _KEY_UP, Name: api.LocalizableString("'Up' move key"), Value: glfw.KeySpace}
	keyDown    = api.TypeKeyValue{Type: reflect.Int, Key: _KEY_DOWN, Name: api.LocalizableString("'Down' move key"), Value: glfw.KeyLeftShift}
	keyExit    = api.TypeKeyValue{Type: reflect.Int, Key: _KEY_EXIT, Name: api.LocalizableString("'Exit' move key"), Value: glfw.KeyEscape}

	_EMPTY_EVENTS = make([]api.Event, 0)
)

type hidGl struct {
	keys   map[glfw.Key]hid.Key
	window *glfw.Window
	closed bool
	time   int64

	Runtime base.Runtime          `inject:""`
	Context modules.GlobalContext `inject:""`
}

func NewMouseKb() modules.HID {
	return &hidGl{
		closed: false,
	}
}

func (this *hidGl) ID() string {
	return "hidGl"
}

func (this *hidGl) Priority() int8 {
	return 127
}

func (this *hidGl) Init() error {
	return this.SetConfiguration(this.Configuration()...)
}

func (this *hidGl) Start() (e error) {
	this.time = time.Now().UnixNano()
	return
}

func (this *hidGl) Stop() {
	this.closed = true
}

func (this *hidGl) Events(t uint64) []api.Event {
	//log.Tracef("%v ms", t)

	ret := _EMPTY_EVENTS

	if this.window == nil {
		log.Debugf("waiting context %#v", this.Context.Get("window"))

		this.window, _ = this.Context.Get("window").(*glfw.Window)

		return ret
	}

	t1 := time.Now().UnixNano()
	dt := t1 - this.time
	
	
	x, y := this.window.GetCursorPos()
	w, h := this.window.GetSize()

	this.window.SetCursorPos(float64(w)*0.5, float64(h)*0.5)

	dh, dv := (float64(w)*0.5)-x, (float64(h)*0.5)-y

	pos := api.Coords{X: 0.0, Y: 0.0, Z: 0.0}

	for k, v := range this.keys {
		if this.window.GetKey(k) == glfw.Press {
			switch v {
			case hid.Forward:
				log.Tracef("Forward move")
				pos.X += 1
			case hid.Back:
				log.Tracef("Back move")
				pos.X -= 1
			case hid.Left:
				log.Tracef("Left move")
				pos.Z -= 1
			case hid.Right:
				log.Tracef("Right move")
				pos.Z += 1
			case hid.Up:
				log.Tracef("Up move")
				pos.Y += 1
			case hid.Down:
				log.Tracef("Down move")
				pos.Y -= 1
			case hid.Exit:
				log.Tracef("Exit HID")
				//this.Stop()
				this.Runtime.Stop()
				this.closed = true
			}
		}
	}

	if dt < 1 {
		return ret
	}
	
	if pos.X != 0 || pos.Y != 0 || pos.Z != 0 || dh != 0 || dv != 0 {
		res := &hid.Action{}
		res.Offset = pos
		res.HorizontalAngle, res.VerticalAngle = dh, dv
		res.AcTime = dt
		log.Debugf("Pos changed: %#v (%v / %v) for %v ns (%v -> %v)", pos, dh, dv, dt, this.time, t1)

		ret = []api.Event{res}
	}

	this.time = t1
	
	return ret
}

func (this *hidGl) Configuration() []api.TypeKeyValue {
	return []api.TypeKeyValue{
		keyForward, keyBack, keyLeft, keyRight, keyUp, keyDown, keyExit,
	}
}

func (this *hidGl) SetConfiguration(values ...api.TypeKeyValue) error {
	this.keys = make(map[glfw.Key]hid.Key)

	var err *multierror.Error

	for _, v := range values {
		key, ok := v.Value.(glfw.Key)

		if !ok {
			err = multierror.Append(err, errors.New(fmt.Sprintf("broken value for %s - %v", v.Key, v.Value)))
			continue
		}

		switch v.Key {
		case _KEY_FORWARD:
			this.keys[key] = hid.Forward
		case _KEY_BACK:
			this.keys[key] = hid.Back
		case _KEY_LEFT:
			this.keys[key] = hid.Left
		case _KEY_RIGHT:
			this.keys[key] = hid.Right
		case _KEY_UP:
			this.keys[key] = hid.Up
		case _KEY_DOWN:
			this.keys[key] = hid.Down
		case _KEY_EXIT:
			this.keys[key] = hid.Exit
		}
	}

	return err.ErrorOrNil()
}
