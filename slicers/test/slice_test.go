package test

import (
	. "github.com/l1va/goosli/primitives"
	"github.com/l1va/goosli/slicers"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestSlicers_SliceByVectorZ(t *testing.T) {
	cases := []struct {
		in     string
		layers int
	}{
		{
			in:     "bokal",
			layers: 800,
		},
		{
			in:     "leyka",
			layers: 500,
		},
		{
			in:     "odn",
			layers: 520,
		},
		{
			in:     "odnzak",
			layers: 575,
		},
		{
			in:     "truba",
			layers: 398,
		},
	}

	for _, row := range cases {
		t.Run(row.in, func(t *testing.T) {
			mesh, err := LoadSTL("../test_models/" + row.in + ".stl")
			if err != nil {
				log.Fatal("failed to load mesh: ", err)
			}
			sett := slicers.Settings{
				LayerHeight:   0.2,
				WallThickness: 1.2,
				FillDensity:   20,
				LineWidth:     0.4,
			}

			gcd := slicers.SliceByVectorToGcode(mesh, AxisZ, sett)

			require.Equal(t, row.layers, gcd.LayersCount)
		})
	}
}
