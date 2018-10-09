package layers

import (
	//"math"
	
	"github.com/iand/perlin"
	"github.com/snuk182/anlgo" 
)

type perlinGenerator struct {
	alpha, beta, scale func(float64) float64
	seeder func(int64) int64
	octaves func(int) int
}

func (this *perlinGenerator) fill(x,z int32, seed int64, a, b, scale float64, octaves int) (float64) {
	xf,zf := float64(x),float64(z)
	return perlin.Noise2D(xf/this.scale(scale), zf/this.scale(scale), this.seeder(seed), this.alpha(a), this.beta(b), this.octaves(octaves))
}

type anlGenerator struct {
	anl.ImplicitModule
}

func (this *anlGenerator) init(seed int64) {
	groundGradient := anl.NewImplicitGradient()
	groundGradient.SetGradient(0,0,0,1,0,0,0,0,0,0,0,0)
	
	//fractaltype=anl.BILLOW, basistype=anl.GRADIENT, interptype=anl.QUINTIC, octaves=2, frequency=0.25}
	lowlandShapeFractal := anl.NewImplicitFractal(anl.Billow, anl.Gradval, anl.Quintic)
	lowlandShapeFractal.SetNumOctaves(2)
	lowlandShapeFractal.SetFrequency(0.25)
	lowlandShapeFractal.SetSeed(uint32(seed))
	
	lowlandAutocorrect := anl.NewImplicitAutoCorrect(0, 1)
	lowlandAutocorrect.SetSource(lowlandShapeFractal)
	
	//{name="lowland_scale",                 type="scaleoffset",      source="lowland_autocorrect", scale=0.125, offset=-0.45},
	lowlandScale := anl.NewImplicitScaleOffset(125, -0.45)
	lowlandScale.SetSourceModule(lowlandAutocorrect)
	
	//{name="lowland_y_scale",               type="scaledomain",      source="lowland_scale", scaley=0},
	lowlandYScale := anl.NewImplicitScaleDomain(1,0,1,1,1,1)
	lowlandYScale.SetSourceModule(lowlandScale)
	
	//{name="lowland_terrain",               type="translatedomain",  source="ground_gradient", ty="lowland_y_scale"},
	lowlandTerrain := anl.NewImplicitTranslateDomain()
	lowlandTerrain.SetSourceModule(groundGradient)
	lowlandTerrain.SetYAxisSourceModule(lowlandYScale)
	
	this.ImplicitModule = lowlandScale
}

func (this *anlGenerator) fill(x,z int32, seed int64, a, b, scale float64, octaves int) (float64) {
	xf,zf := float64(x),float64(z)
	return this.Get2D(xf/scale, zf/scale)
}