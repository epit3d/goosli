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