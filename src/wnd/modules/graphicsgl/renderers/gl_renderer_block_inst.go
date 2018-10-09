package renderers

import (
	"wnd/api"
	//"wnd/base/graphics"
	"wnd/modules/graphicsgl/innerapi"
	"wnd/utils/log"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/jragonmiris/mathgl"
)

const (
	_UVX_BLOCK = 0.3333 / 4.0
	_UVY_BLOCK = 0.5 / 4.0
)

var (
	SKY     = [3]float32{1., 1., 1.}
	_GROUND = [3]float32{0.7, 0.7, 0.7}
)

var verticesBlockNoSlope = []float32{
	//north
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.0, 0.0,
	//south
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.0, 0.9999,
	//east
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	//west
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	//top
	0.0, 0.9999, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.9999, 0.9999,
	//bottom
	0.0, 0.0, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.0, 0.9999,
	0.0, 0.0, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.0, 0.0,
}

var normalsBlockNoSlope = []float32{
	//north
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	//south
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	//east
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	//west
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	//top
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	//bottom
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
}

var uvsNoSlope = []float32{
	//north
	0.0, 0.4999,
	0.0, 0.0,
	0.3332, 0.0,
	0.0, 0.4999,
	0.3332, 0.0,
	0.3332, 0.4999,
	//south
	0.6667, 0.4999,
	0.6667, 0.0,
	1.0, 0.0,
	0.6667, 0.4999,
	1.0, 0.0,
	1.0, 0.4999,
	//east
	0.0, 1.0,
	0.0, 0.5001,
	0.3332, 0.5001,
	0.0, 1.0,
	0.3332, 0.5001,
	0.3332, 1.0,
	//west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
	//top
	0.3334, 0.4999,
	0.3334, 0.0,
	0.6665, 0.0,
	0.3334, 0.4999,
	0.6665, 0.0,
	0.6665, 0.4999,
	//bottom
	0.3334, 1.0,
	0.3334, 0.5001,
	0.6665, 0.5001,
	0.3334, 1.0,
	0.6665, 0.5001,
	0.6665, 1.0,
}

var verticesBlockTopSouth = []float32{
	//south
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.0, 0.9999,
	//east
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	//west
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.0, 0.0,
	//bottom
	0.0, 0.0, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.0, 0.9999,
	0.0, 0.0, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.0, 0.0,
	//top-north
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.0,
}

var normalsBlockTopSouth = []float32{
	//south
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	//east
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	//west
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	//bottom
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	//top-north
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
	0.0, 1.0, -1.0,
}

var uvsBlockTopSouth = []float32{
	//south
	0.6667, 0.4999,
	0.6667, 0.0,
	1.0, 0.0,
	0.6667, 0.4999,
	1.0, 0.0,
	1.0, 0.4999,
	//east
	0.0, 1.0,
	0.3332, 0.5001,
	0.3332, 1.0,
	//west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 1.0,
	//bottom
	0.3334, 1.0,
	0.3334, 0.5001,
	0.6665, 0.5001,
	0.3334, 1.0,
	0.6665, 0.5001,
	0.6665, 1.0,
	//top-north
	0.0, 0.4999,
	0.0, 0.0,
	0.3332, 0.0,
	0.0, 0.4999,
	0.3332, 0.0,
	0.3332, 0.4999,
}

var verticesBlockBottomSouth = []float32{
	//south
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.0, 0.9999,
	//east
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	//west
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	//top
	0.0, 0.9999, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.9999, 0.9999,
	//bottom-north
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.0,
	0.9999, 0.0, 0.9999,
	0.0, 0.9999, 0.0,
	0.0, 0.0, 0.9999,
}

var normalsBlockBottomSouth = []float32{
	//south
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	//east
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	//west
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	//bottom
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	//bottom-north
	0.0, -1.0, -1.0,
	0.0, -1.0, -1.0,
	0.0, -1.0, -1.0,
	0.0, -1.0, -1.0,
	0.0, -1.0, -1.0,
	0.0, -1.0, -1.0,
}

var uvsBlockBottomSouth = []float32{
	//south
	0.6667, 0.4999,
	0.6667, 0.0,
	1.0, 0.0,
	0.6667, 0.4999,
	1.0, 0.0,
	1.0, 0.4999,
	//east
	0.3332, 1.0,
	0.0, 0.5001,
	0.3332, 0.5001,
	//west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	//top
	0.3334, 0.4999,
	0.3334, 0.0,
	0.6665, 0.0,
	0.3334, 0.4999,
	0.6665, 0.0,
	0.6665, 0.4999,
	//bottom-north
	0.0, 0.4999,
	0.0, 0.0,
	0.3332, 0.0,
	0.0, 0.4999,
	0.3332, 0.0,
	0.3332, 0.4999,
}

var verticesBlockTopNorth = []float32{
	//north
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.0, 0.0,
	//east
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.0, 0.9999,
	//west
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	//bottom
	0.0, 0.0, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.0, 0.9999,
	0.0, 0.0, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.0, 0.0,
	//top-south
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.9999,
}

var normalsBlockTopNorth = []float32{
	//north
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	//east
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	//west
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	//bottom
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	//top-south
	0.0, 1.0, 1.0,
	0.0, 1.0, 1.0,
	0.0, 1.0, 1.0,
	0.0, 1.0, 1.0,
	0.0, 1.0, 1.0,
	0.0, 1.0, 1.0,
}

var uvsBlockTopNorth = []float32{
	//north
	0.0, 0.4999,
	0.0, 0.0,
	0.3332, 0.0,
	0.0, 0.4999,
	0.3332, 0.0,
	0.3332, 0.4999,
	//east
	0.0, 1.0,
	0.0, 0.5001,
	0.3332, 1.0,
	//west
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
	//bottom
	0.3334, 1.0,
	0.3334, 0.5001,
	0.6665, 0.5001,
	0.3334, 1.0,
	0.6665, 0.5001,
	0.6665, 1.0,
	//top-south
	0.6667, 0.4999,
	0.6667, 0.0,
	1.0, 0.0,
	0.6667, 0.4999,
	1.0, 0.0,
	1.0, 0.4999,
}

var verticesBlockBottomNorth = []float32{
	//north
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.0, 0.0,
	//east
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	//west
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	//top
	0.0, 0.9999, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.9999, 0.9999,
	//bottom-south
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.0, 0.0, 0.0,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.0, 0.0,
}

var normalsBlockBottomNorth = []float32{
	//north
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	//east
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	//west
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	//top
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	//bottom-south
	0.0, -1.0, 1.0,
	0.0, -1.0, 1.0,
	0.0, -1.0, 1.0,
	0.0, -1.0, 1.0,
	0.0, -1.0, 1.0,
	0.0, -1.0, 1.0,
}

var uvsBlockBottomNorth = []float32{
	//north
	0.0, 0.4999,
	0.0, 0.0,
	0.3332, 0.0,
	0.0, 0.4999,
	0.3332, 0.0,
	0.3332, 0.4999,
	//east
	0.0, 1.0,
	0.0, 0.5001,
	0.3332, 0.5001,
	//west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
	//top
	0.3334, 0.4999,
	0.3334, 0.0,
	0.6665, 0.0,
	0.3334, 0.4999,
	0.6665, 0.0,
	0.6665, 0.4999,
	//bottom-south
	0.6667, 0.4999,
	0.6667, 0.0,
	1.0, 0.0,
	0.6667, 0.4999,
	1.0, 0.0,
	1.0, 0.4999,
}

var verticesBlockTopEast = []float32{
	//north
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.0, 0.0,
	//south
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.0, 0.9999,
	//west
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	//bottom
	0.0, 0.0, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.0, 0.9999,
	0.0, 0.0, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.0, 0.0,
	//top-east
	0.0, 0.0, 0.0,
	0.9999, 0.9999, 0.0,
	0.9999, 0.9999, 0.9999,
	0.0, 0.0, 0.0,
	0.9999, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
}

var normalsBlockTopEast = []float32{
	//north
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	//south
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	//west
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	//bottom
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	//top-east
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
}

var uvsBlockTopEast = []float32{
	//north
	0.0, 0.4999,
	0.0, 0.0,
	0.3332, 0.4999,
	//south
	0.6667, 0.4999,
	1.0, 0.0,
	1.0, 0.4999,
	//west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
	//bottom
	0.3334, 1.0,
	0.3334, 0.5001,
	0.6665, 0.5001,
	0.3334, 1.0,
	0.6665, 0.5001,
	0.6665, 1.0,
	//top-east
	0.0, 1.0,
	0.0, 0.5001,
	0.3332, 0.5001,
	0.0, 1.0,
	0.3332, 0.5001,
	0.3332, 1.0,
}

var verticesBlockTopWest = []float32{
	//north
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.0, 0.0,
	//south
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.0, 0.9999,
	//east
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	//bottom
	0.0, 0.0, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.0, 0.9999,
	0.0, 0.0, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.0, 0.0,
	//top-west
	0.9999, 0.0, 0.9999,
	0.0, 0.9999, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.0, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
}

var normalsBlockTopWest = []float32{
	//north
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	//south
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	//east
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	//bottom
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	//top-west
	1.0, 1.0, 0.0,
	1.0, 1.0, 0.0,
	1.0, 1.0, 0.0,
	1.0, 1.0, 0.0,
	1.0, 1.0, 0.0,
	1.0, 1.0, 0.0,
}

var uvsBlockTopWest = []float32{
	//north
	0.0, 0.4999,
	0.3332, 0.0,
	0.3332, 0.4999,
	//south
	0.6667, 0.4999,
	0.6667, 0.0,
	1.0, 0.4999,
	//east
	0.0, 1.0,
	0.0, 0.5001,
	0.3332, 0.5001,
	0.0, 1.0,
	0.3332, 0.5001,
	0.3332, 1.0,
	//bottom
	0.3334, 1.0,
	0.3334, 0.5001,
	0.6665, 0.5001,
	0.3334, 1.0,
	0.6665, 0.5001,
	0.6665, 1.0,
	//top-west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
}

var verticesBlockBottomEast = []float32{
	//north
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.0,
	//south
	0.9999, 0.0, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.9999,
	//west
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	//top
	0.0, 0.9999, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.9999, 0.9999,
	//bottom-east
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.9999, 0.0, 0.9999,
}

var normalsBlockBottomEast = []float32{
	//north
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	//south
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	//west
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	//top
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	//bottom-east
	-1.0, -1.0, 0.0,
	-1.0, -1.0, 0.0,
	-1.0, -1.0, 0.0,
	-1.0, -1.0, 0.0,
	-1.0, -1.0, 0.0,
	-1.0, -1.0, 0.0,
}

var uvsBlockBottomEast = []float32{
	//north
	0.0, 0.4999,
	0.0, 0.0,
	0.3332, 0.0,
	//south
	1.0, 0.4999,
	0.6667, 0.4999,
	0.6667, 0.0,
	//west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
	//top
	0.3334, 0.4999,
	0.3334, 0.0,
	0.6665, 0.0,
	0.3334, 0.4999,
	0.6665, 0.0,
	0.6665, 0.4999,
	//bottom-east
	0.0, 1.0,
	0.0, 0.5001,
	0.3332, 0.5001,
	0.0, 1.0,
	0.3332, 0.5001,
	0.3332, 1.0,
}

var verticesBlockBottomWest = []float32{
	//north
	0.0, 0.0, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.0,
	//south
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.9999,
	//east
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	//top
	0.0, 0.9999, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.9999, 0.9999,
	//bottom-west
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.0, 0.0, 0.0,
}

var normalsBlockBottomWest = []float32{
	//north
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	//south
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	//east
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	//top
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	//bottom-west
	1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
}

var uvsBlockBottomWest = []float32{
	//north
	0.3332, 0.4999,
	0.0, 0.0,
	0.3332, 0.0,
	//south
	0.6667, 0.4999,
	0.6667, 0.0,
	1.0, 0.0,
	//east
	0.0, 1.0,
	0.0, 0.5001,
	0.3332, 0.5001,
	0.0, 1.0,
	0.3332, 0.5001,
	0.3332, 1.0,
	//top
	0.3334, 0.4999,
	0.3334, 0.0,
	0.6665, 0.0,
	0.3334, 0.4999,
	0.6665, 0.0,
	0.6665, 0.4999,
	//bottom-west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
}

var verticesBlockSouthEast = []float32{
	//south
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.0, 0.9999,
	//west
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	//top
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.9999, 0.9999,
	//bottom
	0.9999, 0.0, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.0, 0.9999,
	//north-east
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
}

var normalsBlockSouthEast = []float32{
	//south
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	//west
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	//top
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	//bottom
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	//north-east
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
	-1.0, 0.0, -1.0,
}

var uvsBlockSouthEast = []float32{
	//south
	0.6667, 0.4999,
	0.6667, 0.0,
	1.0, 0.0,
	0.6667, 0.4999,
	1.0, 0.0,
	1.0, 0.4999,
	//west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
	//top
	0.3334, 0.4999,
	0.6665, 0.0,
	0.6665, 0.4999,
	//bottom
	0.6665, 1.0,
	0.3334, 0.5001,
	0.6665, 0.5001,
	//north-east
	0.0, 0.4999,
	0.0, 0.0,
	0.3332, 0.0,
	0.0, 0.4999,
	0.3332, 0.0,
	0.3332, 0.4999,
}

var verticesBlockSouthWest = []float32{
	//south
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.0, 0.9999,
	//east
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	//top
	0.0, 0.9999, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.9999,
	//bottom
	0.0, 0.0, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.0, 0.9999,
	//north-west
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.0, 0.9999,
	0.0, 0.9999, 0.0,
	0.0, 0.0, 0.0,
}

var normalsBlockSouthWest = []float32{
	//south
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	//east
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	//top
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	//bottom
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	//north-west
	1.0, 0.0, -1.0,
	1.0, 0.0, -1.0,
	1.0, 0.0, -1.0,
	1.0, 0.0, -1.0,
	1.0, 0.0, -1.0,
	1.0, 0.0, -1.0,
}

var uvsBlockSouthWest = []float32{
	//south
	0.6667, 0.4999,
	0.6667, 0.0,
	1.0, 0.0,
	0.6667, 0.4999,
	1.0, 0.0,
	1.0, 0.4999,
	//east
	0.0, 1.0,
	0.0, 0.5001,
	0.3332, 0.5001,
	0.0, 1.0,
	0.3332, 0.5001,
	0.3332, 1.0,
	//top
	0.3334, 0.4999,
	0.3334, 0.0,
	0.6665, 0.4999,
	//bottom
	0.3334, 1.0,
	0.3334, 0.5001,
	0.6665, 0.5001,
	//north-west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
}

var verticesBlockNorthEast = []float32{
	//north
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.0, 0.0,
	//west
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	//top
	0.9999, 0.9999, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.0,
	//bottom
	0.0, 0.0, 0.0,
	0.9999, 0.0, 0.9999,
	0.9999, 0.0, 0.0,
	//south-east
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.9999,
	0.0, 0.0, 0.0,
	0.9999, 0.9999, 0.9999,
	0.9999, 0.0, 0.9999,
}

var normalsBlockNorthEast = []float32{
	//north
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	//west
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	//top
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	//bottom
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	//south-east
	-1.0, 0.0, 1.0,
	-1.0, 0.0, 1.0,
	-1.0, 0.0, 1.0,
	-1.0, 0.0, 1.0,
	-1.0, 0.0, 1.0,
	-1.0, 0.0, 1.0,
}

var uvsBlockNorthEast = []float32{
	//north
	0.0, 0.4999,
	0.0, 0.0,
	0.3332, 0.0,
	0.0, 0.4999,
	0.3332, 0.0,
	0.3332, 0.4999,
	//west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
	//top
	0.6665, 0.4999,
	0.3334, 0.0,
	0.6665, 0.0,
	//bottom
	0.3334, 1.0,
	0.6665, 0.5001,
	0.6665, 1.0,
	//south-east
	0.6667, 0.4999,
	0.6667, 0.0,
	1.0, 0.0,
	0.6667, 0.4999,
	1.0, 0.0,
	1.0, 0.4999,
}

var verticesBlockNorthWest = []float32{
	//north
	0.9999, 0.0, 0.0,
	0.9999, 0.9999, 0.0,
	0.0, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.0, 0.0,
	//east
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.0,
	0.0, 0.9999, 0.9999,
	0.0, 0.0, 0.9999,
	//top
	0.0, 0.9999, 0.9999,
	0.0, 0.9999, 0.0,
	0.9999, 0.9999, 0.0,
	//bottom
	0.0, 0.0, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.0, 0.0,
	//south-west
	0.0, 0.0, 0.9999,
	0.0, 0.9999, 0.9999,
	0.9999, 0.9999, 0.0,
	0.0, 0.0, 0.9999,
	0.9999, 0.9999, 0.0,
	0.9999, 0.0, 0.0,
}

var normalsBlockNorthWest = []float32{
	//north
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	0.0, 0.0, -1.0,
	//east
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	-1.0, 0.0, 0.0,
	//top
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
	//bottom
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	0.0, -1.0, 0.0,
	//south-west
	1.0, 0.0, 1.0,
	1.0, 0.0, 1.0,
	1.0, 0.0, 1.0,
	1.0, 0.0, 1.0,
	1.0, 0.0, 1.0,
	1.0, 0.0, 1.0,
}

var uvsBlockNorthWest = []float32{
	//north
	0.0, 0.4999,
	0.0, 0.0,
	0.3332, 0.0,
	0.0, 0.4999,
	0.3332, 0.0,
	0.3332, 0.4999,
	//east
	0.0, 1.0,
	0.0, 0.5001,
	0.3332, 0.5001,
	0.0, 1.0,
	0.3332, 0.5001,
	0.3332, 1.0,
	//top
	0.3334, 0.4999,
	0.3334, 0.0,
	0.6665, 0.0,
	//bottom
	0.3334, 1.0,
	0.3334, 0.5001,
	0.6665, 1.0,
	//south-west
	0.6667, 1.0,
	0.6667, 0.5001,
	1.0, 0.5001,
	0.6667, 1.0,
	1.0, 0.5001,
	1.0, 1.0,
}

type bufferObject struct {
	data []float32
	id uint32
}

func newBufferObject() *bufferObject {
	return &bufferObject{
		data: make([]float32, 0),
	}
}

type blockRenderData struct {
	id        innerapi.RenderDataID
	programId innerapi.RenderProgramID

	primitiveCount int
	instanceCount  int
	
	vertices, textureUVs, normals, positions bufferObject
	
	vao uint32
	shift mathgl.Vec3f

	p innerapi.RenderProgram
}

func newBlockRenderData(id innerapi.RenderDataID, pid innerapi.RenderProgramID, size api.BlockSize, slope api.BlockSlope) *blockRenderData {
	log.Tracef("graphicsgl.newBlockRenderData: id: %s, pid: %s", id, pid)

	rrr := new(blockRenderData)
	rrr.id = id
	rrr.programId = pid
	
	switch slope {
	case api.SlopeNone:
		rrr.vertices.data = verticesBlockNoSlope
		rrr.textureUVs.data = uvsNoSlope
		rrr.normals.data = normalsBlockNoSlope
	case api.SlopeEastBottom:
		rrr.vertices.data = verticesBlockBottomEast
		rrr.textureUVs.data = uvsBlockBottomEast
		rrr.normals.data = normalsBlockBottomEast
	case api.SlopeEastTop:
		rrr.vertices.data = verticesBlockTopEast
		rrr.textureUVs.data = uvsBlockTopEast
		rrr.normals.data = normalsBlockTopEast
	case api.SlopeSouthBottom:
		rrr.vertices.data = verticesBlockBottomSouth
		rrr.textureUVs.data = uvsBlockBottomSouth
		rrr.normals.data = normalsBlockBottomSouth
	case api.SlopeSouthEast:
		rrr.vertices.data = verticesBlockSouthEast
		rrr.textureUVs.data = uvsBlockSouthEast
		rrr.normals.data = normalsBlockSouthEast
	case api.SlopeSouthTop:
		rrr.vertices.data = verticesBlockTopSouth
		rrr.textureUVs.data = uvsBlockTopSouth
		rrr.normals.data = normalsBlockTopSouth
	case api.SlopeSouthWest:
		rrr.vertices.data = verticesBlockSouthWest
		rrr.textureUVs.data = uvsBlockSouthWest
		rrr.normals.data = normalsBlockSouthWest
	case api.SlopeNorthBottom:
		rrr.vertices.data = verticesBlockBottomNorth
		rrr.textureUVs.data = uvsBlockBottomNorth
		rrr.normals.data = normalsBlockBottomNorth
	case api.SlopeNorthEast:
		rrr.vertices.data = verticesBlockNorthEast
		rrr.textureUVs.data = uvsBlockNorthEast
		rrr.normals.data = normalsBlockNorthEast
	case api.SlopeNorthTop:
		rrr.vertices.data = verticesBlockTopNorth
		rrr.textureUVs.data = uvsBlockTopNorth
		rrr.normals.data = normalsBlockTopNorth
	case api.SlopeNorthWest:
		rrr.vertices.data = verticesBlockNorthWest
		rrr.textureUVs.data = uvsBlockNorthWest
		rrr.normals.data = normalsBlockNorthWest
	case api.SlopeWestBottom:
		rrr.vertices.data = verticesBlockBottomWest
		rrr.textureUVs.data = uvsBlockBottomWest
		rrr.normals.data = normalsBlockBottomWest
	case api.SlopeWestTop:
		rrr.vertices.data = verticesBlockTopWest
		rrr.textureUVs.data = uvsBlockTopWest
		rrr.normals.data = normalsBlockTopWest
	default:
		rrr.vertices.data = []float32{}
		rrr.textureUVs.data = []float32{}
		rrr.normals.data = []float32{}
	}

	/*for i:=0; i<len(rrr.instanceVertices); i+=3 {

			}*/

	rrr.primitiveCount = len(rrr.vertices.data) / 3

	if size != api.SizeFull {
		multiplier := size.Multiplier()

		tmpv := make([]float32, len(rrr.vertices.data))
		copy(tmpv, rrr.vertices.data)

		for i, c := range tmpv {
			tmpv[i] = c * multiplier
		}

		rrr.vertices.data = tmpv

		tmpu := make([]float32, len(rrr.textureUVs.data))
		copy(tmpu, rrr.textureUVs.data)

		for i := 0; i < len(tmpu); i += 2 {
			//tmpu[i+1] = tmpu[i] + (_UVY_BLOCK)
		}

		rrr.textureUVs.data = tmpu
	}

	return rrr
}

func (this *blockRenderData) Clone() *blockRenderData {
	return &blockRenderData{
		id:        this.id,
		programId: this.programId,

		primitiveCount: this.primitiveCount,
		instanceCount:  this.instanceCount,
		
		vertices: this.vertices,
		textureUVs: this.textureUVs,
		normals: this.normals,

		positions:          this.positions,

		p: this.p,
	}
}

func (this *blockRenderData) RenderDataID() innerapi.RenderDataID {
	return this.id
}

func (this *blockRenderData) RenderProgramIDs() []innerapi.RenderProgramID {
	return []innerapi.RenderProgramID{this.programId}
}

func (this *blockRenderData) ID() string {
	return string(this.id)
}

func (this *blockRenderData) Target() string {
	return "graphics/gl"
}

func (this *blockRenderData) Init(p []innerapi.RenderProgram) {
	if this.vertices.data != nil {
		gl.DeleteVertexArrays(1, &this.vao)
		gl.GenVertexArrays(1, &this.vao)

		gl.DeleteBuffers(1, &this.vertices.id)
		gl.GenBuffers(1, &this.vertices.id)
		gl.BindBuffer(gl.ARRAY_BUFFER, this.vertices.id)
		gl.BufferData(gl.ARRAY_BUFFER, len(this.vertices.data)*4, gl.Ptr(this.vertices.data), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		this.vertices.data = nil

		gl.DeleteBuffers(1, &this.textureUVs.id)
		gl.GenBuffers(1, &this.textureUVs.id)
		gl.BindBuffer(gl.ARRAY_BUFFER, this.textureUVs.id)
		gl.BufferData(gl.ARRAY_BUFFER, len(this.textureUVs.data)*4, gl.Ptr(this.textureUVs.data), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		this.textureUVs.data = nil

		gl.DeleteBuffers(1, &this.normals.id)
		gl.GenBuffers(1, &this.normals.id)
		gl.BindBuffer(gl.ARRAY_BUFFER, this.normals.id)
		gl.BufferData(gl.ARRAY_BUFFER, len(this.normals.data)*4, gl.Ptr(this.normals.data), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		this.textureUVs.data = nil

		if len(p) > 0 {
			this.p = p[0]
		} else {
			log.Errorf("No RenderProgramID provided for %v", this.programId)
		}
	}
	
	if this.positions.data != nil {
		gl.DeleteBuffers(1, &this.positions.id)
		gl.GenBuffers(1, &this.positions.id)
		gl.BindBuffer(gl.ARRAY_BUFFER, this.positions.id)
		gl.BufferData(gl.ARRAY_BUFFER, len(this.positions.data)*4, gl.Ptr(this.positions.data), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}
}

var ld = mathgl.Vec3f{0.0, 1.0, 0.0}

func (this *blockRenderData) Render(camera *innerapi.GlCamera) {
	gl.Enable(gl.CULL_FACE)
	gl.FrontFace(gl.CW)

	gl.BindVertexArray(this.vao)

	gl.UseProgram(this.p.Program)

	gl.UniformMatrix4fv(this.p.MVP, 1, false, &camera.MVP[0])
	gl.UniformMatrix4fv(this.p.M, 1, false, &camera.M[0])
	gl.UniformMatrix4fv(this.p.V, 1, false, &camera.V[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, this.p.Texture)
	gl.Uniform1i(this.p.Texs, 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, this.p.NormalTexture)
	gl.Uniform1i(this.p.Norm, 1)

	gl.Uniform3fv(this.p.LightDir, 1, &ld[0])
	gl.Uniform3fv(this.p.SkyColor, 1, &SKY[0])
	gl.Uniform3fv(this.p.GroundColor, 1, &_GROUND[0])
	
	gl.Uniform3fv(this.p.Shift, 1, &this.shift[0])

	gl.EnableVertexAttribArray(this.p.Vertex)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vertices.id)
	gl.VertexAttribPointer(this.p.Vertex, 3, gl.FLOAT, false, 0, gl.Ptr(nil))

	gl.EnableVertexAttribArray(this.p.Color)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.textureUVs.id)
	gl.VertexAttribPointer(this.p.Color, 2, gl.FLOAT, false, 0, gl.Ptr(nil))

	gl.EnableVertexAttribArray(this.p.Normal)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.normals.id)
	gl.VertexAttribPointer(this.p.Normal, 3, gl.FLOAT, false, 0, gl.Ptr(nil))

	gl.EnableVertexAttribArray(this.p.Position)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.positions.id)
	gl.VertexAttribPointer(this.p.Position, 3, gl.FLOAT, false, 0, gl.Ptr(nil))
	gl.VertexAttribDivisor(this.p.Position, 1)

	gl.DrawArraysInstanced(gl.TRIANGLES, 0, int32(this.primitiveCount), int32(this.instanceCount))

	gl.DisableVertexAttribArray(this.p.Vertex)
	gl.DisableVertexAttribArray(this.p.Color)
	gl.DisableVertexAttribArray(this.p.Position)
	gl.DisableVertexAttribArray(this.p.Normal)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	gl.UseProgram(0)

	gl.BindVertexArray(0)
}

func (this *blockRenderData) Purge() {
	gl.DeleteVertexArrays(1, &this.vao)
	gl.DeleteBuffers(1, &this.vertices.id)
	gl.DeleteBuffers(1, &this.textureUVs.id)
	gl.DeleteBuffers(1, &this.normals.id)
	gl.DeleteBuffers(1, &this.positions.id)
}
