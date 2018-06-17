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
		out Line
	}{
		{
			in1: Plane{Point{0, 0, 0}, V(1, 1, 1)},
			in2: NewTriangle(Point{0, 0, -1}, Point{1, 0, -1}, Point{0, 0, 1}),
			out: Line{P1: Point{1,0,-1}, P2:Point{0,0,0}},
		},
		{
			in1: Plane{Point{0, 0, 0}, V(0, 0, 1)},
			in2: NewTriangle(Point{0, 0, -1}, Point{2, 0, -1}, Point{0, 0, 1}),
			out: Line{P1: Point{1,0,0}, P2:Point{0,0,0}},
		},
		{
			in1: Plane{Point{0, 0, 0}, V(0, 0, 1)},
			in2: NewTriangle(Point{0, 0, -1}, Point{2, 0, -1}, Point{0, 0, 1}),
			out: Line{P1: Point{1,0,0}, P2:Point{0,0,0}},
		},
	}
	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			res := row.in1.IntersectTriangle(&row.in2)
			require.Equal(t, row.out, *res)
		})
	}
}
