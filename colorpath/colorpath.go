package colorpath

import (
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/lucasb-eyer/go-colorful"
)

// implements the map object interface for go-staticmaps for a path object that can
// vary the color of the path along the way

// Point is a coordinate and a color
type Point struct {
	s2.LatLng
	Color colorful.Color
}

// ColorPath satisfies the map object interface for go-staticmap
type ColorPath struct {
	sm.MapObject
	Positions []Point
	Weight    float64
}

// NewColorPath builds a new path with colors
func NewColorPath(weight float64) *ColorPath {
	cp := new(ColorPath)
	cp.Positions = []Point{}
	cp.Weight = weight
	return cp
}

// ExtraMarginPixels - to help go-staticmap find render bounds
// its just a line so no padding
func (cp *ColorPath) ExtraMarginPixels() (float64, float64, float64, float64) {
	return cp.Weight, cp.Weight, cp.Weight, cp.Weight
}

// Bounds returns the geographical boundary rect (excluding the actual pixel dimensions).
func (cp *ColorPath) Bounds() s2.Rect {
	r := s2.EmptyRect()
	for _, ll := range cp.Positions {
		r = r.AddPoint(ll.LatLng)
	}
	return r
}

// Draw draws the colorpath in the given graphical context.
func (cp *ColorPath) Draw(gc *gg.Context, trans *sm.Transformer) {
	if len(cp.Positions) <= 1 {
		return
	}

	gc.ClearPath()
	gc.SetLineWidth(cp.Weight)
	gc.SetLineCap(gg.LineCapRound)
	gc.SetLineJoin(gg.LineJoinRound)

	for i := 1; i < len(cp.Positions); i++ {
		gc.SetColor(cp.Positions[i-1].Color)
		spx, spy := trans.LatLngToXY(cp.Positions[i-1].LatLng)
		epx, epy := trans.LatLngToXY(cp.Positions[i].LatLng)
		gc.DrawLine(spx, spy, epx, epy)
		if i%2 == 0 {
			gc.Stroke()
		}
	}
	gc.Stroke()
}
