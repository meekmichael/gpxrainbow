package positionregistry

import (
	"math"
	"sync"

	"github.com/golang/geo/s2"
	"github.com/meekmichael/gpxrainbow/colorpath"
)

// PositionRegistry keeps state for prior locations seen in GPX paths
type PositionRegistry struct {
	MaxColors uint16
	SeenPos   map[int][]s2.LatLng
	Tracks    int
}

// CountNear returns the number of previous paths that had one point within _meters_
// meters of the current position.  This approach, while kind of slow, yields a
// good result because we don't double count us backtracking over the same
// point on any one walk, and doesn't count when the points in a segment are less
// than _meters_ meters apart.
func (p *PositionRegistry) CountNear(ll s2.LatLng, meters float64) uint16 {
	ret := uint16(0)
	wg := sync.WaitGroup{}
	retChan := make(chan uint16, len(p.SeenPos))
	for _, positions := range p.SeenPos {
		wg.Add(1)

		// for each prior path, concurrently check if we've been close to this spot
		go func(wg *sync.WaitGroup, ll s2.LatLng, meters float64, positions []s2.LatLng, rChan chan uint16) {
			defer wg.Done()
			for _, prevPt := range positions {
				rad := prevPt.Distance(ll).Radians()         // degrees apart that these two points arex
				distance := 40_050_000 / (2 * math.Pi) * rad // approx circumference of earth, close enough for a toy rainbow map program
				if distance < meters {
					rChan <- 1
					return
				}
			}
			rChan <- 0
		}(&wg, ll, meters, positions, retChan)

	}
	wg.Wait()

	for range p.SeenPos {
		ret += <-retChan
	}

	return ret
}

// AddFromColorPath adds all the positions in a colorpath to the registry
func (p *PositionRegistry) AddFromColorPath(cp *colorpath.ColorPath, trkNum int) {
	if p.SeenPos == nil {
		p.SeenPos = map[int][]s2.LatLng{}
	}
	for _, pos := range cp.Positions {
		p.SeenPos[trkNum] = append(p.SeenPos[trkNum], pos.LatLng)
	}
}
