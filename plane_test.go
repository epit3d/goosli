package goosli

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/require"
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
