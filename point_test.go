package goosli

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/require"
)

func TestPoint_Equal(t *testing.T) {
	cases := []struct {
		in1, in2 Point
		out      bool
	}{
		{
			in1: Point{1, 0, 0},
			in2: Point{0, 1, 0},
			out: false,
		},
		{
			in1: Point{1, -1, 0},
			in2: Point{1, -1, 0},
			out: true,
		},
		{
			in1: Point{0, 0, 0},
			in2: Point{0, 0, 0},
			out: true,
		},
		{
			in1: Point{33.333333, 1, 1},
			in2: Point{33.3333333333333, 1, 1},
			out: false,
		},
		{
			in1: Point{33.333333333, 1, 1},
			in2: Point{33.3333333333333, 1, 1},
			out: true,
		},
		{
			in1: Point{33.3333333333333, 1, 1},
			in2: Point{33.333333333, 1, 1},
			out: true,
		},
		{
			in1: Point{33.3333333333333, 1, 1},
			in2: Point{33.333333, 1, 1},
			out: false,
		},
	}

	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			res := row.in1.Equal(row.in2)
			require.Equal(t, row.out, res)
		})
	}
}
