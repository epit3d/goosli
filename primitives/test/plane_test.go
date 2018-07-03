package test

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/require"
	. "github.com/l1va/goosli/primitives"
)

func TestPlane_Intersect(t *testing.T) {
	cases := []struct {
		in1 Plane
		in2 Triangle
		out bool
	}{
		{
			in1: Plane{Point{1, 1, 1}, V(1, 1, 1)},
			in2: NewTriangle(Point{1, 0, 1}, Point{1, 2, 1}, Point{0, 1, 1}),
			out: true,
		},
		{
			in1: Plane{Point{0, 0, 0}, V(0, 0, 1)},
			in2: NewTriangle(Point{1, 1, 1}, Point{2, 2, 2}, Point{0, 1, 1}),
			out: false,
		},
		{
			in1: Plane{Point{-1.072904954318182, 19.34564047988636, 28.694603830000002}, V(0, 0, 1)},
			in2: NewTriangle( //Point{-0.14450382574562326 ,-0.028103494667281348 ,0.004398115494216481},
				Point{-2.999999999999999, 19.92450142132615, 28.695752504523085},
				Point{-2.99210286, 17.22907075536739, 11.731701764504587},
				Point{-2.99044621, 17.220552517303418, 11.731701765622354}),
			out: true,
		},
	}
	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			res := row.in1.Intersect(&row.in2)
			require.Equal(t, row.out, res)
		})
	}
}

func TestPlane_IntersectTriangle(t *testing.T) {
	cases := []struct {
		in1 Plane
		in2 Triangle
		out *Line
	}{
		{
			in1: Plane{Point{0, 0, 0}, V(1, 1, 1)},
			in2: NewTriangle(Point{0, 0, -1}, Point{1, 0, -1}, Point{0, 0, 1}),
			out: &Line{P1: Point{1, 0, -1}, P2: Point{0, 0, 0}},
		},
		{
			in1: Plane{Point{0, 0, 0}, V(0, 0, 1)},
			in2: NewTriangle(Point{0, 0, -1}, Point{2, 0, -1}, Point{0, 0, 1}),
			out: &Line{P1: Point{1, 0, 0}, P2: Point{0, 0, 0}},
		},
		{
			in1: Plane{Point{0, 0, 0}, V(0, 0, 1)},
			in2: NewTriangle(Point{0, 0, -1}, Point{2, 0, -1}, Point{0, 0, 1}),
			out: &Line{P1: Point{1, 0, 0}, P2: Point{0, 0, 0}},
		},
		{

			in1: Plane{Point{-4.290768404625, -8.1802342555, 20.891189639999997}, V(0, 0, 1)},
			in2: NewTriangle(Point{38.563079833984375, -38.947452545166016, 29.953892707824707},
				Point{33.27024841308594, -34.81993103027344, 19.940381050109863},
				Point{33.028175354003906, -31.475749969482422, 20.79648494720459}),
			out: &Line{P1: Point{33.0854166, -31.55302132, 20.89118964}, P2: Point{33.77281632, -35.21184977, 20.89118964}},
		},
		{
			in1: Plane{Point{0, 0, 0}, V(0, 0, 1)},
			in2: NewTriangle(Point{0, 0, 0}, Point{2, 0, 1}, Point{0, 0, 1}),
			out: nil,
		},
		{
			in1: Plane{Point{0, 0, 0}, V(0, 0, 1)},
			in2: NewTriangle(Point{0, 0, -0.000000001}, Point{2, 0, 1}, Point{0, 0, 1}),
			out: nil,
		},
		{
			in1: Plane{Point{0, 0, 0}, V(0, 0, 1)},
			in2: NewTriangle(Point{0, 0, -0.00001}, Point{2, 0, 1}, Point{0, 0, 1}),
			out: &Line{P1: Point{X: 2e-05, Y: 0, Z: 0}, P2: Point{X: 0, Y: 0, Z: 0}},
		},
		{//TODO: fixme
			in1: Plane{Point{-1.072904954318182, 19.34564047988636, 28.694603830000002}, V(0, 0, 1)},
			in2: NewTriangle( //Point{-0.14450382574562326 ,-0.028103494667281348 ,0.004398115494216481},
				Point{-2.999999999999999, 19.92450142132615, 28.695752504523085},
				Point{-2.99210286, 17.22907075536739, 11.731701764504587},
				Point{-2.99044621, 17.220552517303418, 11.731701765622354}),
			out: &Line{P1: Point{-2.9999994652666584, 19.92431890757503, 28.694603830000002}, P2: Point{-2.9999993530911127, 19.924318330785713, 28.694603830000002}},
		},
	}
	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			res := row.in1.IntersectTriangle(&row.in2)
			if row.out == nil {
				require.Nil(t, res)
			} else {
				require.Equal(t, *row.out, *res)
			}
		})
	}
}
