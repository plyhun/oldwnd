package api

type LocalizableString string

type GameObject interface {
	ID() string
}

type Observer struct {
	Position, LooksAt Coords
}

func (this *Observer) Direction() []Direction {
	dx, dy, dz := this.LooksAt.X - this.Position.X, this.LooksAt.Y - this.Position.Y, this.LooksAt.Z - this.Position.Z
	
	out := make([]Direction, 0, 3)
	
	if dx != 0 {
		if dx > 0 {
			out = append(out, DirectionWest)
		} else {
			out = append(out, DirectionEast)
		}
	}
	
	if dz != 0 {
		if dz > 0 {
			out = append(out, DirectionNorth)
		} else {
			out = append(out, DirectionSouth)
		}
	}
	
	if dy != 0 {
		if dy > 0 {
			out = append(out, DirectionUp)
		} else {
			out = append(out, DirectionDown)
		}
	}
	
	return out
}

//Adjacency map _NSEWTBD N - north, S - south, E - east, W - west, T - top, B- bottom, D - definition bit (0 - undefined, 1 - defined)

const (
	bitExtra      uint8 = 0x80
	bitNorth      uint8 = 0x40
	bitSouth      uint8 = 0x20
	bitEast       uint8 = 0x10
	bitWest       uint8 = 0x8
	bitTop        uint8 = 0x4
	bitBottom     uint8 = 0x2
	bitDefinition uint8 = 0x1
)

type AdjacencyList uint8

func (this AdjacencyList) Add(dir ...Direction) AdjacencyList {
	for _, d := range dir {
		switch d {
		case DirectionNorth:
			this |= AdjacencyList(bitNorth)
		case DirectionSouth:
			this |= AdjacencyList(bitSouth)
		case DirectionEast:
			this |= AdjacencyList(bitEast)
		case DirectionWest:
			this |= AdjacencyList(bitWest)
		case DirectionUp:
			this |= AdjacencyList(bitTop)
		case DirectionDown:
			this |= AdjacencyList(bitBottom)
		}
	}
	
	this = this.SetDefined(true)
	
	return this
}

func (this AdjacencyList) Remove(dir ...Direction) AdjacencyList {
	for _, d := range dir {
		switch d {
		case DirectionNorth:
			this &= AdjacencyList(^bitNorth)
		case DirectionSouth:
			this &= AdjacencyList(^bitSouth)
		case DirectionEast:
			this &= AdjacencyList(^bitEast)
		case DirectionWest:
			this &= AdjacencyList(^bitWest)
		case DirectionUp:
			this &= AdjacencyList(^bitTop)
		case DirectionDown:
			this &= AdjacencyList(^bitBottom)
		}
	}
	
	this = this.SetDefined(true)
	
	return this
}

func (this AdjacencyList) SetDefined(defined bool) AdjacencyList {
	if defined {
		this |= AdjacencyList(bitDefinition)
	} else {
		this &= AdjacencyList(^bitDefinition)
	}
	
	return this
}

func (this AdjacencyList) IsDefined() bool {
	return (uint8(this) & bitDefinition) > 0
}

func (this AdjacencyList) HasAllAdjacents() bool {
	return (uint8(this) & 0x7f) >= 0x7f
}

func (this AdjacencyList) HasNoAdjacents() bool {
	return (uint8(this) & 0x7f) <= 1
}

func (this AdjacencyList) Inverse() AdjacencyList {
	ret := ^this
	ret = ret.SetDefined(this.IsDefined())
	
	return ret
} 

func (this AdjacencyList) Parse() []Direction {
	ret := make([]Direction, 6)
	
	count := 0
	
	if (uint8(this) & bitNorth) > 0 {
		ret[count] = DirectionNorth
		count++
	}
	
	if (uint8(this) & bitSouth) > 0 {
		ret[count] = DirectionSouth
		count++
	}
	
	if (uint8(this) & bitEast) > 0 {
		ret[count] = DirectionEast
		count++
	}
	
	if (uint8(this) & bitWest) > 0 {
		ret[count] = DirectionWest
		count++
	}
	
	if (uint8(this) & bitTop) > 0 {
		ret[count] = DirectionUp
		count++
	}
	
	if (uint8(this) & bitBottom) > 0 {
		ret[count] = DirectionDown
		count++
	}
	
	return ret[:count]
} 