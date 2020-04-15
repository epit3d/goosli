package test

import (
	"fmt"
	//"fmt"
	"testing"
	//"github.com/stretchr/testify/require"
	. "github.com/l1va/goosli/primitives"
	"github.com/stretchr/testify/require"
)

func TestPath_Join(t *testing.T) { //TODO: add more and improve
	var cases = []struct {
		in  []Path
		out []Path
	}{
		{
			in: []Path{
				{Points: []Point{{1, 1, 1}, {1, 1, 2}}},
				{Points: []Point{{1, 1, 3}, {1, 1, 2}}},
				{Points: []Point{{1, 1, 3}, {1, 1, 4}}},
				{Points: []Point{{1, 1, 1}, {1, 1, 4}}},
			},
			out: []Path{{Points: []Point{
				{1, 1, 4}, {1, 1, 1}, {1, 1, 2}, {1, 1, 3}, {1, 1, 4}}}},
		},
		{
			in: []Path{
				{Points: []Point{{1, 1, 1}, {1, 2, 1}}},
				{Points: []Point{{1, 4, 1}, {1, 1, 1}}},
				{Points: []Point{{1, 3, 1}, {1, 4, 1}}},
				{Points: []Point{{1, 2, 1}, {1, 3, 1}}},
			},
			out: []Path{{Points: []Point{{1, 3, 1}, {1, 4, 1}, {1, 1, 1}, {1, 2, 1},
				{1, 3, 1},}}},
		},
		{
			in: []Path{
				{Points: []Point{{1, 1, 1}, {1, 2, 1}, {1, 3, 1}}},
				{Points: []Point{{1, 1, 1}, {1, 2, 1}}},
				{Points: []Point{{1, 1, 1}, {1, 2, 1}, {1, 3, 1}}},
				{Points: []Point{{1, 2, 1}, {1, 3, 1}}},
			},
			out: []Path{{Points: []Point{{1, 1, 1}, {1, 2, 1}, {1, 3, 1},
				{1, 2, 1}, {1, 1, 1}, {1, 2, 1}, {1, 3, 1}}}},
		},
		{
			in: []Path{
				{Points: []Point{{1, 1, 1}, {2, 1, 1}, {3, 1, 1}}},
				{Points: []Point{{4, 1, 1}, {5, 1, 1}}},
			},
			out: []Path{
				{Points: []Point{{1, 1, 1}, {2, 1, 1}, {3, 1, 1}}},
				{Points: []Point{{4, 1, 1}, {5, 1, 1}}},},
		},
		{
			in: []Path{
				{Points: []Point{{-19.840, -5.845, 69.351}, {-20.7470000001, -5.844999999, 67.408}}},
				{Points: []Point{{-20.747, -5.845, 67.408}, {-22.116, -5.845, 64.472}}},
				{Points: []Point{{-24.630, -5.845, 31.285}, {-24.673, -5.845, 0.000}}},
			},
			out: []Path{
				{Points: []Point{{-24.630, -5.845, 31.285}, {-24.673, -5.845, 0.000}}},
				{Points: []Point{{-19.840, -5.845, 69.351}, {-20.747, -5.845, 67.408}, {-22.116, -5.845, 64.472}}},
			},
		},
		{
			in: []Path{
				{Points: []Point{{-19.240, -6.690, 50.000}, {-19.150, -6.820, 51.640}}},
				{Points: []Point{{-18.970, -7.060, 56.430}, {-19.150, -6.820, 51.640}}},
				{Points: []Point{{-19.240, -6.690, 50.000}, {-19.230, -6.700, 40.190}}},
			},
			out: []Path{
				{Points: []Point{{-18.970, -7.060, 56.430}, {-19.150, -6.820, 51.640}, {-19.240, -6.690, 50.000}, {-19.230, -6.700, 40.190}}},
			},
		},
		{
			in: []Path{
				{Points: []Point{{-22.050, -5.850, 64.610}, {-24.650, -5.850, 50.000}}},
				{Points: []Point{{-24.630, -5.850, 30.950}, {-24.650, -5.850, 50.000}}},
			},
			out: []Path{
				{Points: []Point{{-22.050, -5.850, 64.610}, {-24.650, -5.850, 50.000}, {-24.630, -5.850, 30.950}}},
			},
		},
	}
	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			res := JoinPaths3(row.in)
			require.Equal(t, len(row.out), len(res))
			for _, p := range res {
				found := false
				for _, p2 := range row.out {
					if p.Equal(p2) {
						found = true
					}
				}
				require.True(t, found, "not found: ", p)
			}
		})
	}
}

func TestPath_FindCentroid(t *testing.T) {
	cases := []struct {
		in  Path
		out Point
	}{
		{
			in:  Path{Points: []Point{{1, 1, 1}, {1, 1, 2}, {1, 2, 2}, {X: 1, Y: 2, Z: 1}, {1, 1, 1}}},
			out: Point{X: 1, Y: 1.5, Z: 1.5},
		},
		{
			in:  Path{Points: []Point{{2, 1, 1}, {2, 1, 2}, {1, 2, 2}, {1, 2, 1}, {2, 1, 1}}},
			out: Point{1.5, 1.5, 1.5},
		},
		{
			in:  Path{Points: []Point{{2, 1, 1}, {2, 2, 2}, {1, 3, 1}, {1, 2, 3}, {2, 1, 1}},},
			out: Point{1.5, 2, 2},
		},
		{
			in:  Path{Points: []Point{{15, 3, 2}, {3, 2, 12}, {1, 13, 10}, {12, 11, 3}, Point{15, 3, 2}},},
			out: Point{X: 7.472222222222222, Y: 7.138888888888889, Z: 7.444444444444445},
		},
	}
	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require.Equal(t, row.out, FindCentroid(row.in))
		})
	}
}
