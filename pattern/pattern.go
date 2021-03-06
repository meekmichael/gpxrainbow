package pattern

import (
	"log"

	"github.com/lucasb-eyer/go-colorful"
)

type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

func MustHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// GetInterpolatedColorFor is borrowed from https://github.com/lucasb-eyer/go-colorful
func GetGradientTable() GradientTable {
	// I'm red/green colorblind and I picked these colors because they're easy
	// for me to see.  YMMV.  TODO: make this customizable

	// TODO greyscale, reds, blues, greens
	return GradientTable{
		{MustHex("#001eff"), 0},
		{MustHex("#2db2ee"), 0.2},
		{MustHex("#00ff78"), 0.4},
		{MustHex("#deef03"), 0.6},
		{MustHex("#b18d1e"), 0.8},
		{MustHex("#c92009"), 1},
	}
}

func (gt GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt[len(gt)-1].Col
}
