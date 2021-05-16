package legend

import (
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/meekmichael/gpxrainbow/pattern"
)

const outerXPt = 330
const outerYPt = 120
const outerYHeight = 105
const rainbowWidth = 280
const rainbowYTop = 90
const rainbowHeight = 50

var outerColor = color.RGBA{0, 0, 0, 180}
var textColor = color.RGBA{248, 248, 248, 255}

// Options configures the legend
type Options struct {
	FormatString  string
	GradientTable pattern.GradientTable
	MaxVal        float64
	MinVal        float64
	Steps         int
	Title         string
}

// Render puts the legend on the image and returns the be-legened image
func Render(opts Options, img image.Image) (image.Image, error) {
	gc := gg.NewContextForImage(img)
	if gc.Width() < (2*outerXPt) || gc.Height() < (3*outerYHeight) {
		// image too small for legend
		return img, nil
	}
	if opts.Steps < 2 {
		return img, nil
	}
	gc.SetColor(outerColor)
	gc.DrawRectangle(float64(gc.Width())-outerXPt, float64(gc.Height())-outerYPt, outerXPt, outerYHeight)
	gc.Fill()

	for i := 0; i < opts.Steps; i++ {
		step := float64(i) * float64(rainbowWidth/opts.Steps)
		gc.DrawRectangle(float64(gc.Width())-rainbowWidth-20+step, float64(gc.Height())-rainbowYTop, rainbowWidth-step, rainbowHeight)
		gc.SetColor(pattern.GetGradientTable().GetInterpolatedColorFor(float64(i) / float64(opts.Steps)))
		gc.Fill()
	}
	gc.SetColor(textColor)
	gc.DrawStringAnchored(opts.Title, float64(gc.Width())-(rainbowWidth/2)-20, float64(gc.Height()-rainbowYTop+rainbowHeight+10), 0.5, 0.5)
	lSteps := 4.0
	if float64(opts.Steps) < lSteps {
		lSteps = float64(opts.Steps)
	}
	lStep := 0.0
	for i := opts.MinVal; i <= opts.MaxVal; i += ((opts.MaxVal - opts.MinVal) / lSteps) {
		gc.DrawStringAnchored(fmt.Sprintf(opts.FormatString, i), float64(gc.Width())-rainbowWidth-20+(rainbowWidth*lStep), float64(gc.Height())-outerYHeight+5, 0.5, 0.5)
		lStep += (1.0 / lSteps)
	}

	return gc.Image(), nil
}
