package path

import (
	"errors"
	"fmt"
	"math"
	"strings"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/meekmichael/gpxrainbow/colorpath"
	"github.com/meekmichael/gpxrainbow/config"
	"github.com/meekmichael/gpxrainbow/legend"
	"github.com/meekmichael/gpxrainbow/pattern"
	"github.com/meekmichael/gpxrainbow/positionregistry"
	"github.com/meekmichael/gpxrainbow/tile"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/urfave/cli/v2"
)

// Run is the main method for this project
func Run(c *cli.Context) error {
	ctx := sm.NewContext()
	mConf, err := config.NewConfig(c)
	if err != nil {
		return err
	}
	ctx.SetSize(mConf.ImageWidth, mConf.ImageHeight)
	ctx.SetTileProvider(tile.ProviderByName(mConf.TileProvider))

	gpxFiles := c.Args().Slice()

	if len(gpxFiles) == 0 {
		return errors.New("no file(s) specified")
	}
	if mConf.Mode != config.MODE_PROXIMITY {
		pathData, err := maxSpeedAndElev(gpxFiles)
		if err != nil {
			return err
		}
		mConf.MinElevation = pathData.MinElevation
		mConf.MaxElevation = pathData.MaxElevation
		mConf.MaxSpeed = pathData.MaxSpeed
	}

	paths := []*colorpath.ColorPath{}
	fmt.Printf("Processing %d GPX files\n", len(gpxFiles))
	posRegistry := positionregistry.PositionRegistry{
		MaxColors: uint16(len(gpxFiles)),
	}
	for _, gpxFile := range gpxFiles {
		newPaths, err := gpxToColorPath(mConf, gpxFile, &posRegistry)
		if err != nil {
			return err
		}
		paths = append(paths, newPaths...)
	}
	for _, p := range paths {
		ctx.AddObject(p)
	}
	img, err := ctx.Render()

	if err != nil {
		return err
	}
	legendOpts := legend.Options{
		GradientTable: pattern.GetGradientTable(),
	}
	switch mConf.Mode {
	case config.MODE_PROXIMITY:
		legendOpts.MinVal = 1
		legendOpts.MaxVal = float64(posRegistry.Tracks)
		legendOpts.FormatString = "%2.0f"
		legendOpts.Steps = posRegistry.Tracks
		legendOpts.Title = "count"
	case config.MODE_ELEVATION:
		legendOpts.Steps = 250
		if mConf.Units == "us" {
			legendOpts.MinVal = mConf.MinElevation * 3.2808399
			legendOpts.MaxVal = mConf.MaxElevation * 3.2808399
			legendOpts.Title = "elevation (ft)"
		} else {
			legendOpts.MinVal = mConf.MinElevation
			legendOpts.MaxVal = mConf.MaxElevation
			legendOpts.Title = "elevation (meters)"
		}
		legendOpts.FormatString = "%2.0f"
	case config.MODE_SPEED:
		legendOpts.MinVal = 0
		if mConf.Units == "us" {
			legendOpts.MaxVal = mConf.MaxSpeed * 2.236936 // meters/s -> mph\
			legendOpts.Title = "Speed (mph)"
		} else {
			legendOpts.MaxVal = mConf.MaxSpeed * 3.6 // meters/s -> kph
			legendOpts.Title = "speed (kph)"
		}
		legendOpts.Steps = 250

		legendOpts.FormatString = "%2.1f"
	}

	if mConf.Mode != config.MODE_INPUT {
		img, err = legend.Render(legendOpts, img)
	}
	if err != nil {
		return err
	}

	if strings.HasSuffix(mConf.OutputFile, "png") || strings.HasSuffix(mConf.OutputFile, "PNG") {
		if err := gg.SavePNG(mConf.OutputFile, img); err != nil {
			return err
		}
	} else {
		if err := gg.SaveJPG(mConf.OutputFile, img, 85); err != nil {
			return err
		}
	}
	if err == nil {
		fmt.Printf("Saved as %s\n", mConf.OutputFile)
	}
	return err
}

// AggregatePathData is aggregate information about all paths togethe
type AggregatePathData struct {
	MinElevation float64
	MaxElevation float64
	MaxSpeed     float64
}

func maxSpeedAndElev(filenames []string) (AggregatePathData, error) {
	minElev := math.Inf(1)
	maxElev := 0.0
	maxSpeed := 0.0
	for _, gpxFile := range filenames {
		gpxdata, err := gpx.ParseFile(gpxFile)
		if err != nil {
			return AggregatePathData{}, fmt.Errorf("likely invalid GPX file %s, error: %v", gpxFile, err)
		}
		for _, trk := range gpxdata.Tracks {
			for _, seg := range trk.Segments {
				if err != nil {
					return AggregatePathData{}, err
				}
				for i := range seg.Points {
					elev := seg.Points[i].GetElevation()
					if elev.NotNull() {
						if elev.Value() > maxElev {
							maxElev = elev.Value()
						}
						if elev.Value() < minElev {
							minElev = elev.Value()
						}
					}
					if i == 0 {
						continue
					}
					spd := seg.Points[i].SpeedBetween(&seg.Points[i-1], true)
					if spd > maxSpeed && spd != math.Inf(1) {
						maxSpeed = spd
					}
				}
			}
		}
	}
	return AggregatePathData{
		MinElevation: minElev,
		MaxElevation: maxElev,
		MaxSpeed:     maxSpeed,
	}, nil
}

// gpxToColorPath iterates through a single GPX file and builds a ColorPath object
// to be later drawn onto a map
func gpxToColorPath(conf config.MapConfig, filename string, posRegistry *positionregistry.PositionRegistry) ([]*colorpath.ColorPath, error) {
	gpxdata, err := gpx.ParseFile(filename)
	if err != nil {
		return []*colorpath.ColorPath{}, err
	}
	paths := []*colorpath.ColorPath{}
	elevDiff := conf.MaxElevation - conf.MinElevation
	for _, trk := range gpxdata.Tracks {
		if conf.Mode == config.MODE_INPUT || conf.Mode == config.MODE_PROXIMITY {
			posRegistry.Tracks++
		}
		for _, seg := range trk.Segments {
			lastColor := pattern.GetGradientTable().GetInterpolatedColorFor(0)
			p := colorpath.NewColorPath(float64(conf.LineWidth))
			spd := float64(0)
			for i := 0; i < len(seg.Points); i++ {
				color := colorful.Color{}
				if i > 0 {
					spd = seg.Points[i].SpeedBetween(&seg.Points[i-1], true) // meters/second
				}
				elev := seg.Points[i].Elevation
				switch conf.Mode {
				case config.MODE_INPUT:
					color = pattern.GetGradientTable().GetInterpolatedColorFor(float64(posRegistry.Tracks) / float64(posRegistry.MaxColors))
				case config.MODE_PROXIMITY:
					countNear := uint16(0)
					if posRegistry.Tracks > 1 {
						countNear = posRegistry.CountNear(s2.LatLngFromDegrees(seg.Points[i].GetLatitude(), seg.Points[i].GetLongitude()), float64(conf.ProximityDistance))
					}
					color = pattern.GetGradientTable().GetInterpolatedColorFor(float64(countNear) / float64(posRegistry.MaxColors))
				case config.MODE_SPEED:
					// blend factor to make segments blend together better and not be wild colors
					// especially useful for when there are a lot of points close together in a segment
					color = pattern.GetGradientTable().GetInterpolatedColorFor(spd/conf.MaxSpeed).BlendHcl(lastColor, 0.7)
				case config.MODE_ELEVATION:
					if elev.NotNull() {
						color = pattern.GetGradientTable().GetInterpolatedColorFor((elev.Value()-conf.MinElevation)/elevDiff).BlendHcl(lastColor, 0.5)
					} else {
						color = lastColor
					}
				}
				lastColor = color
				p.Positions = append(p.Positions, colorpath.Point{
					Color:  color,
					LatLng: s2.LatLngFromDegrees(seg.Points[i].GetLatitude(), seg.Points[i].GetLongitude()),
				})
			}
			paths = append(paths, p)
			if conf.Mode == config.MODE_PROXIMITY {
				posRegistry.AddFromColorPath(p, posRegistry.Tracks)
			}
			if conf.Mode == config.MODE_INPUT || conf.Mode == config.MODE_PROXIMITY {
				posRegistry.Tracks++
			}
		}
	}
	return paths, nil
}
