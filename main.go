package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/meekmichael/gpxrainbow/config"
	"github.com/meekmichael/gpxrainbow/path"
	"github.com/urfave/cli/v2"
)

func main() {
	logger := log.New(os.Stdout, "gpxrainbow: ", log.Lshortfile)
	app := &cli.App{
		Name:     "gpxrainbow",
		HelpName: "",
		Usage:    "plots GPX tracks on an openstreetmap map with the color of the path conveying additional information",
		Version:  "1.0",
		Commands: nil,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "width",
				Aliases: []string{"x"},
				Usage:   "width of output image",
				Value:   2048,
			},
			&cli.IntFlag{
				Name:    "height",
				Aliases: []string{"y"},
				Usage:   "height of output image",
				Value:   1536,
			},
			&cli.IntFlag{
				Name:    "linewidth",
				Aliases: []string{"l"},
				Usage:   "line width (in pixels)",
				Value:   3,
			},
			&cli.StringFlag{
				Name:    "mode",
				Aliases: []string{"m"},
				Usage:   "mode - [proximity|speed|elevation|date]",
				Value:   config.MODE_PROXIMITY,
			},
			&cli.StringFlag{
				Name:    "tileprovider",
				Aliases: []string{"tp"},
				Usage:   "OpenStreetMap tile provider, use --list-tileprovider to get a list",
				Value:   "carto-light",
			},
			&cli.BoolFlag{
				Name:  "list-tileprovider",
				Usage: "list available tileproviders for --tileprovider",
				Value: false,
			},
			&cli.StringFlag{
				Name:    "outputfile",
				Aliases: []string{"o"},
				Usage:   "file to write the map to (must be a png)",
				Value:   "output.png",
			},
			&cli.IntFlag{
				Name:    "proximity_distance",
				Aliases: []string{"d"},
				Usage:   "distance in meters (approx) to color path the same in proximity mode",
				Value:   10,
			},
			&cli.StringFlag{
				Name:    "units",
				Aliases: []string{"u"},
				Usage:   "units - \"us\" or \"metric\"",
				Value:   "metric",
			},
		},
		Action: path.Run,

		OnUsageError: func(*cli.Context, error, bool) error {
			return errors.New("invalid arguments (try --help)")
		},
		Compiled: time.Time{},
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
