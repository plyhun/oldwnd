package entity

import (
	"wnd/api"
	"wnd/api/events"
	"wnd/utils/log"
	"wnd/base/graphics"
	
	"math"
	
	glm "github.com/go-gl/mathgl/mgl32"
)

const (
	_HALF_PI = math.Pi / 2.0
	_HALF_PI_LESS = _HALF_PI - 0.1
)

type Playable struct {
	*api.Entity
	va, ha float64
}

func NewPlayable(e *api.Entity) *Playable {
	return &Playable{Entity: e, va: 0.0, ha: 0.0}
}

func (this *Playable) MakeMove(event *events.EntityMove) *events.EntityPosition {
	if event == nil || this.Entity == nil {
		return nil
	}

	log.Tracef("%v -> %v", this.Entity.ID, event)
	
	pos := glm.Vec3{this.Position.X, this.Position.Y, this.Position.Z}
	//lookAt := glm.Vec3f{this.Entity.LooksAt.X, this.Entity.LooksAt.Y, this.Entity.LooksAt.Z}

	accel := float32(event.DeltaTime) / 100000000000 * float32(this.Entity.Speed)

	this.va += event.VerticalAngle * float64(accel) / 2
	this.ha += event.HorizontalAngle * float64(accel) / 2
	
	if this.va > _HALF_PI_LESS {
		this.va = _HALF_PI_LESS
	} else if this.va < -_HALF_PI_LESS {
		this.va = -_HALF_PI_LESS
	} 

	if this.ha > 360 {
		this.ha -= 360
	} else if this.ha < -360 {
		this.ha += 360
	} 

	fb := glm.Vec3{
		float32(math.Cos(this.va) * math.Sin(this.ha)),
		float32(math.Sin(this.va)),
		float32(math.Cos(this.va) * math.Cos(this.ha)),
	}
	lr := glm.Vec3{
		float32(math.Sin(this.ha - _HALF_PI)),
		0.0,
		float32(math.Cos(this.ha - _HALF_PI)),
	}
	
	//log.Tracef("fb %#v\nlr %#v\nud %#v", fb, lr, graphics.UD)

	accel *= float32(this.Entity.Speed)
	
	if event.Offset.X > 0 {
		dfb := fb.Mul(accel)
		dfb[1] = 0
		pos = pos.Add(dfb)
	}
	if event.Offset.X < 0 {
		dfb := fb.Mul(accel)
		dfb[1] = 0
		pos = pos.Sub(dfb)
	}
	if event.Offset.Z > 0 {
		dlr := lr.Mul(accel)
		dlr[1] = 0
		pos = pos.Add(dlr)
	}
	if event.Offset.Z < 0 {
		dlr := lr.Mul(accel)
		dlr[1] = 0
		pos = pos.Sub(dlr)
	}
	if event.Offset.Y > 0 {
		pos = pos.Add(graphics.UD.Mul(accel))
	}
	if event.Offset.Y < 0 {
		pos = pos.Sub(graphics.UD.Mul(accel))
	}

	//log.Tracef("pos moved to %#v, look at %#v", pos, pos.Add(fb))
	
	lookAt := pos.Add(fb)
	
	e := &events.EntityPosition{
		General: event.General,
		EntityID: event.EntityID,
		Observer: api.Observer{
			Position: api.Coords{X: pos[0], Y: pos[1], Z: pos[2]},
			LooksAt: api.Coords{X: lookAt[0], Y: lookAt[1], Z: lookAt[2]},
		},					
	}
	
	//log.Warnf("delta %v / %v / %v (%v / %v)", e.Position.X - this.Position.X, e.Position.Y - this.Position.Y, e.Position.Z - this.Position.Z, this.ha, this.va)

	return e
}
