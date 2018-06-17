package test

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/require"
	. "github.com/l1va/goosli/primitives"
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
