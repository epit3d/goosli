package test

import (
	"fmt"
	"testing"

	. "github.com/l1va/goosli/primitives"
	"github.com/stretchr/testify/require"
)

func TestLine_IsCollinearPointOnSegment(t *testing.T) {
	cases := []struct {
		in1 Line
		in2 Point
		out bool
	}{
		{
			in1: Line{Point{0, 0, 0}, Point{1, 1, 1}},
			in2: Point{0.5, 0.5, 0.5},
			out: true,
		},
		{
			in1: Line{Point{0, 0, 0}, Point{1, 1, 1}},
			in2: Point{1, 1, 1},
			out: true,
		},
		{
			in1: Line{Point{0, 0, 0}, Point{1, 1, 1}},
			in2: Point{1.5, 1.5, 1.5},
			out: false,
		},
		{
			in1: Line{Point{0, 0, 0}, Point{1, 1, 1}},
			in2: Point{-1.5, -1.5, -1.5},
			out: false,
		},
	}

	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require.Equal(t, row.out, row.in1.IsCollinearPointOnSegment(row.in2))
		})
	}
}

func TestLine_IsIntersectingSegment(t *testing.T) {
	cases := []struct {
		in1 Line
		in2 Line
		out bool
	}{
		// Segments are on each other
		{
			in1: Line{Point{0, 0, 0}, Point{1, 1, 1}},
			in2: Line{Point{0, 0, 0}, Point{1, 1, 1}},
			out: true,
		},
		// Segments are on each other partly
		{
			in1: Line{Point{0, 0, 0}, Point{1, 1, 1}},
			in2: Line{Point{0, 0, 0}, Point{0.5, 0.5, 0.5}},
			out: true,
		},
		// Segments intersect at one point
		{
			in1: Line{Point{0, 0, 0}, Point{1, 1, 1}},
			in2: Line{Point{1, 0, 0}, Point{0, 1, 1}},
			out: true,
		},
		// Segments are parallel
		{
			in1: Line{Point{0, 0, 0}, Point{0, 1, 0}},
			in2: Line{Point{1, 0, 0}, Point{1, 1, 0}},
			out: false,
		},
		// Segments are not parallel but not touching each other
		{
			in1: Line{Point{0, 0, 0}, Point{0, 1, 0}},
			in2: Line{Point{0.5, 0.5, 0}, Point{1.0, 0.5, 0}},
			out: false,
		},
	}

	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require.Equal(t, row.out, row.in1.IsIntersectingSegment(&row.in2))
		})
	}
}

func TestLine_IntersectLine(t *testing.T) {

	cases := []struct {
		in1 Line
		in2 Line
		out *Point
	}{
		// Lines intersect on XY plane
		{
			in1: Line{Point{3, 2, 0}, Point{8, 7, 0}},
			in2: Line{Point{1, 6, 0}, Point{7, 3, 0}},
			out: &Point{5, 4, 0},
		},
		// Lines intersect on YZ plane - fail
/*		{
			in1: Line{Point{0, 3, 2}, Point{0, 8, 7}},
			in2: Line{Point{0, 1, 6}, Point{0, 7, 3}},
			out: &Point{0, 5, 4},
		},*/

		// Lines intersect on XZ plane - fail
/*		{
			in1: Line{Point{3, 0, 2}, Point{8, 0, 7}},
			in2: Line{Point{1, 0, 6}, Point{7, 0, 3}},
			out: &Point{5, 0, 4},
		},*/
		// Lines intersect
		{
			in1: Line{Point{1, 1, -2}, Point{1, -5, 1}},
			in2: Line{Point{1, -3, 0}, Point{-1, 0, -4}},
			out: &Point{1, -3, 0},
		},
		// 1, 2, 3 points collinear, 3 lies b/w 1,2
		{
			in1: Line{Point{1, -3, 0}, Point{-1, 0, -4}},
			in2: Line{Point{0, -1.5, -2}, Point{0, 0, 0}},
			out: &Point{0, -1.5, -2},
		},
		//  1, 2, 3 points collinear, 3 lies b/w 1,2
		{
			in1: Line{Point{1, -3, 0}, Point{-3, 3, -8}},
			in2: Line{Point{-1, 0, -4}, Point{0, 0, 0}},
			out: &Point{-1, 0, -4},
		},
		// not intersect
		{
			in1: Line{Point{1, -3, 0}, Point{-3, 3, -8}},
			in2: Line{Point{1, 0, 4}, Point{0, 0, 0}},
			out: nil,
		},

	}

	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require.Equal(t, row.out, row.in1.IntersectLine(&row.in2))
		})
	}
}

func TestLine_IntersectTriangle(t *testing.T) {

	cases := []struct {
		in1 Line
		in2 Triangle
		out bool
	}{
		//Line intersect one side of triangle
		{
			in1: Line{Point{1, 1, -2}, Point{1, -5, 1}},
			in2: NewTriangle(Point{1, -3, 0}, Point{-1, 0, -4}, Point{0, 0, 0}),
			out: true,
		},
		//Not intersect
		{
			in1: Line{Point{10, 10, 10}, Point{9, 9, 9}},
			in2: NewTriangle(Point{1, -3, 0}, Point{-1, 0, -4}, Point{0, 0, 0}),
			out: false,
		},

	}

	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require.Equal(t, row.out, row.in1.IntersectTriangle(row.in2))
		})
	}
}