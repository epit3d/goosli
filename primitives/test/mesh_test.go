package test

import (
	"fmt"
	. "github.com/l1va/goosli/primitives"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMesh_Shift(t *testing.T) {
	cases := []struct {
		in, out Mesh
		in2     Vector
	}{
		{
			in:  NewMesh([]Triangle{NewTriangle(Point{1, 1, 1}, Point{2, 1, 1}, Point{1, 2, 1})}),
			in2: Vector{3, 2, 1},
			out: NewMesh([]Triangle{NewTriangle(Point{4, 3, 2}, Point{5, 3, 2}, Point{4, 4, 2})}),
		},
		{
			in:  NewMesh([]Triangle{NewTriangle(Point{1, 1, 1}, Point{2, 1, 1}, Point{1, 2, 1})}),
			in2: Vector{3, 2, -1},
			out: NewMesh([]Triangle{NewTriangle(Point{4, 3, 0}, Point{5, 3, 0}, Point{4, 4, 0})}),
		},
	}

	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			row.in.Shift(row.in2)
			require.Equal(t, row.out, row.in)
		})
	}
}
