package test

import (
	"fmt"
	. "github.com/l1va/goosli/primitives"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTriangle_MinZ(t *testing.T) {
	cases := []struct {
		in1 Triangle
		in2 Vector
		out float64
	}{
		{
			in1: NewTriangle(Point{1, 0, 1}, Point{1, 2, 1}, Point{0, 1, 1}),
			in2: V(1, 1, 1),
			out: 2,
		},
		{
			in1: NewTriangle(Point{1, 1, 1}, Point{2, 2, 2}, Point{0, 1, 1}),
			in2: V(0, 0, 1),
			out: 1,
		},
	}
	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require.Equal(t, row.out, row.in1.MinZ(row.in2))
		})
	}
}

func TestTriangle_IntersectTriangle(t *testing.T) {
		cases := []struct {
		in1 Triangle
		in2 Triangle
		out *Line
	}{
		//line inside of triangle
		{
			in1: NewTriangle(Point{0, 0, 5}, Point{0, 10, 5}, Point{10, 0, 5}),
			in2: NewTriangle(Point{2, 0, 0}, Point{2, 7, 0}, Point{2, 0, 7}),
			out: &Line{Point{2,0,5}, Point{2,2,5}},
		},
		//line intersect one side of triangle
		{
			in1: NewTriangle(Point{0, 0, 5}, Point{0, 10, 5}, Point{10, 0, 5}),
			in2: NewTriangle(Point{2, -1, 0}, Point{2, 7, 0}, Point{2, -1, 8}),
			out: &Line{Point{2,0,5}, Point{2,2,5}},
		},
		{
			in1: NewTriangle(Point{0, 10, 5}, Point{0, 0, 5}, Point{10, 0, 5}),
			in2: NewTriangle(Point{2, -1, 0}, Point{2, 7, 0}, Point{2, -1, 8}),
			out: &Line{Point{2,2,5}, Point{2,0,5}},
		},
		{
			in1: NewTriangle(Point{0, 0, 5}, Point{10, 0, 5}, Point{0, 10, 5}),
			in2: NewTriangle(Point{2, -1, 0}, Point{2, 7, 0}, Point{2, -1, 8}),
			out: &Line{Point{2,2,5}, Point{2,0,5}},
		},
		//line intersect two sides of triangle
		{
			in1: NewTriangle(Point{0, 0, 5}, Point{10, 0, 5}, Point{0, 10, 5}),
			in2: NewTriangle(Point{9, -1, 0}, Point{9, 7, 0}, Point{9, -1, 8}),
			out: &Line{Point{9,0,5}, Point{9,1,5}},
		},
		{
			in1: NewTriangle(Point{0, 0, 5}, Point{0, 10, 5}, Point{10, 0, 5}),
			in2: NewTriangle(Point{9, -1, 0}, Point{9, 7, 0}, Point{9, -1, 8}),
			out: &Line{Point{9,0,5}, Point{9,1,5}},
		},
		{
			in1: NewTriangle(Point{10, 0, 5}, Point{0, 0, 5}, Point{0, 10, 5}),
			in2: NewTriangle(Point{9, -1, 0}, Point{9, 7, 0}, Point{9, -1, 8}),
			out: &Line{Point{9,0,5}, Point{9,1,5}},
		},
		//line intersect one side and one vertex
		{
			in1: NewTriangle(Point{8, 0, 5}, Point{10, 0, 5}, Point{9, 1, 5}),
			in2: NewTriangle(Point{9, -1, 0}, Point{9, 7, 0}, Point{9, -1, 8}),
			out: &Line{Point{9,0,5}, Point{9,1,5}},
		},
		// not intersect
		{
			in1: NewTriangle(Point{2, 0, 5}, Point{0, 0, 5}, Point{0, 2, 5}),
			in2: NewTriangle(Point{9, -1, 0}, Point{9, 7, 0}, Point{9, -1, 8}),
			out: nil,
		},
	}
	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require.Equal(t, row.out, row.in1.IntersectTriangle(&row.in2))
		})
	}
}

func TestTriangle_PointBelongs(t *testing.T) {
	cases := []struct {
		in1 Triangle
		in2 Point
		out bool
	}{
		//on XY plane
		{
			in1: NewTriangle(Point{10, 0, 5}, Point{0, 0, 5}, Point{0, 10, 5}),
			in2: Point{3,4,5},
			out: true,
		},
		{
			in1: NewTriangle(Point{10, 0, 5}, Point{0, 0, 5}, Point{0, 10, 5}),
			in2: Point{3,4,0},
			out: false,
		},
		//on YZ plane
		{
			in1: NewTriangle(Point{2, 0, 0}, Point{2, 7, 0}, Point{2, 0, 7}),
			in2: Point{2,1,1},
			out: true,
		},
		//on XZ plane
		{
			in1: NewTriangle(Point{0, 2, 0}, Point{7, 2, 0}, Point{0, 2, 7}),
			in2: Point{1,2,1},
			out: true,
		},
		//arbitrary plane
		{
			in1: NewTriangle(Point{2, 2, 2}, Point{4, 1, 1}, Point{6, 2, -2}),
			in2: Point{3,2,1},
			out: true,
		},
		{
			in1: NewTriangle(Point{2, 2, 2}, Point{4, 1, 1}, Point{6, 2, -2}),
			in2: Point{10,-10,6},
			out: false,
		},
	}
	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require.Equal(t, row.out, row.in1.PointBelongs(row.in2))
		})
	}
}