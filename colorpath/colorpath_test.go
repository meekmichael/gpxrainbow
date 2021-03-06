package colorpath

import (
	"image/color"
	"testing"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/meekmichael/gpxrainbow/tile"
	"github.com/stretchr/testify/assert"
)

func TestColorPath_Draw(t *testing.T) {
	type fields struct {
		MapObject sm.MapObject
		Positions []Point
		Weight    float64
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "happypath",
			fields: fields{
				MapObject: &ColorPath{},
				Positions: []Point{
					{
						LatLng: s2.LatLngFromDegrees(45, 45),
						Color: colorful.Color{
							R: 1.0,
							G: 1.0,
							B: 0.0,
						},
					},
					{
						LatLng: s2.LatLngFromDegrees(45.005, 44.995),
						Color: colorful.Color{
							R: 1.0,
							G: 0.95,
							B: 0.05,
						},
					},
					{
						LatLng: s2.LatLngFromDegrees(45.01, 44.99),
						Color: colorful.Color{
							R: 1.0,
							G: 0.9,
							B: 0.1,
						},
					},
				},
				Weight: 3.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := &ColorPath{
				MapObject: tt.fields.MapObject,
				Positions: tt.fields.Positions,
				Weight:    tt.fields.Weight,
			}
			ctx := sm.NewContext()
			ctx.SetSize(64, 64)
			ctx.SetTileProvider(tile.ProviderByName("carto-light"))
			ctx.AddObject(cp)
			// cp.Draw is called
			img, err := ctx.Render()
			assert.NoError(t, err)
			assert.Equal(t, color.RGBA{255, 243, 12, 255}, img.At(32, 32))
			// gg.SavePNG("test.png", img)
		})
	}
}
