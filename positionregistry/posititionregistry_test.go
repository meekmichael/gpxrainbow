package positionregistry

import (
	"testing"

	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

func TestPositionRegistry_CountNear(t *testing.T) {
	type fields struct {
		MaxColors uint16
		SeenPos   map[int][]s2.LatLng
		Tracks    int
	}
	type args struct {
		ll     s2.LatLng
		meters float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint16
	}{
		{
			name: "it runs",
			fields: fields{
				MaxColors: 250,
				SeenPos: map[int][]s2.LatLng{
					1: {
						s2.LatLngFromDegrees(45, 45),
						s2.LatLngFromDegrees(45.001, 45.002),
						s2.LatLngFromDegrees(45.002, 45.002),
						s2.LatLngFromDegrees(45.001, 45.001),
						s2.LatLngFromDegrees(45, 45),
					},
					2: {
						s2.LatLngFromDegrees(45, 45),
						s2.LatLngFromDegrees(45.001, 45.002),
						s2.LatLngFromDegrees(45.005, 45.004),
						s2.LatLngFromDegrees(45.001, 44.997),
						s2.LatLngFromDegrees(45, 45),
					},
				},
			},
			args: args{
				ll:     s2.LatLngFromDegrees(45.005, 45.004),
				meters: 100,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PositionRegistry{
				MaxColors: tt.fields.MaxColors,
				SeenPos:   tt.fields.SeenPos,
				Tracks:    tt.fields.Tracks,
			}
			got := p.CountNear(tt.args.ll, tt.args.meters)
			assert.Equal(t, tt.want, got)
		})
	}
}
