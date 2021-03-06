package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/meekmichael/gpxrainbow/tile"
	"github.com/urfave/cli/v2"
)

// MapConfig is global configuration state
type MapConfig struct {
	ImageHeight       int
	ImageWidth        int
	LineWidth         uint16
	Mode              string
	OutputFile        string
	ProximityDistance uint16
	TileProvider      string
	Units             string

	// set at runtime
	MaxElevation float64
	MaxSpeed     float64
	MinElevation float64
}

// MODE_PROXIMITY color path based on number of proximity to this pixel
const MODE_PROXIMITY = "proximity"

// MODE_INPUT color path based on the order files were read in
const MODE_INPUT = "input"

// MODE_SPEED color path by speed
const MODE_SPEED = "speed"

// MODE_ELEVATION color path by elevation
const MODE_ELEVATION = "elevation"

const maxheight = 16 * 1024
const maxwidth = 16 * 1024
const minheight = 64
const minwidth = 64
const minlinewidth = 1
const maxlinewidth = 32
const minproximity = 1
const maxproximity = 1000

// NewConfig validates inputs and builds a config struct
func NewConfig(c *cli.Context) (MapConfig, error) {
	if c.Bool("list-tileprovider") {
		tile.ListTileProvider()
	}
	tp := c.String("tileprovider")
	if !tile.ValidateTileProvider(tp) {
		return MapConfig{}, fmt.Errorf("invalid tileprovider, use --list-tileprovider to get a list")
	}
	height := c.Int("height")
	if height < minheight || height > maxheight {
		return MapConfig{}, fmt.Errorf("Please use a height between %d and %d", minheight, maxheight)
	}
	width := c.Int("width")
	if width < minwidth || width > maxwidth {
		return MapConfig{}, fmt.Errorf("Please use a width between %d and %d", minwidth, maxwidth)
	}
	lineWidth := c.Int("linewidth")
	if lineWidth < minlinewidth || lineWidth > maxlinewidth {
		return MapConfig{}, fmt.Errorf("Please use a line width between %d and %d", minlinewidth, maxlinewidth)
	}
	proxDistance := c.Int("proximity_distance")
	if proxDistance < minproximity || proxDistance > maxproximity {
		return MapConfig{}, fmt.Errorf("Please use a proximity_distance between %d and %d", minproximity, maxproximity)
	}
	mode := c.String("mode")
	if _, ok := map[string]bool{MODE_PROXIMITY: true, MODE_INPUT: true, MODE_SPEED: true, MODE_ELEVATION: true}[strings.ToLower(mode)]; !ok {
		return MapConfig{}, errors.New("Please pick a valid mode, one of proximity, input, speed, elevation")
	}
	units := strings.ToLower(c.String("units"))
	if units != "us" && units != "metric" {
		return MapConfig{}, errors.New("units must be \"us\" or \"metric\"")
	}

	outfile := filepath.Clean(c.String("outputfile"))

	if !strings.HasSuffix(outfile, ".png") &&
		!strings.HasSuffix(outfile, ".jpg") &&
		!strings.HasSuffix(outfile, ".jpeg") &&
		!strings.HasSuffix(outfile, ".PNG") &&
		!strings.HasSuffix(outfile, ".JPG") &&
		!strings.HasSuffix(outfile, ".JPEG") {
		return MapConfig{}, errors.New("only png and jpeg are supported for the output image")
	}

	return MapConfig{
		ImageHeight:       height,
		ImageWidth:        width,
		LineWidth:         uint16(lineWidth),
		Mode:              mode,
		OutputFile:        outfile,
		ProximityDistance: uint16(proxDistance),
		TileProvider:      tp,
		Units:             units,
	}, nil
}
