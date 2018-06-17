package test

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/require"
	. "github.com/l1va/goosli/primitives"
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


func TestPoint_ProjectOnLine(t *testing.T) {
	cases := []struct {
		in, in1, in2, out Point
	}{
		{
			in: Point{1, 1, 0},
			in1: Point{0, 0, 0},
			in2: Point{0, 1, 0},
			out: Point{0, 1, 0},
		},
		{
			in: Point{1, 1, 0},
			in1: Point{0, 0, 0},
			in2: Point{0, 2, 0},
			out: Point{0, 1, 0},
		},
		{
			in: Point{0, 1, 1},
			in1: Point{0, 0, 0},
			in2: Point{0, 1, 0},
			out: Point{0, 1, 0},
		},
		{
			in: Point{0, 1, 1},
			in1: Point{0, 0, 0},
			in2: Point{0, 2, 0},
			out: Point{0, 1, 0},
		},
		{
			in: Point{1, 1, 0},
			in1: Point{1, 1, 0},
			in2: Point{0, 0, 0},
			out: Point{1, 1, 0},
		},
		{
			in: Point{0, 0, 0},
			in1: Point{-2, 0, 0},
			in2: Point{0, -2, 0},
			out: Point{-1, -1, 0},
		},
	}

	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			res := row.in.ProjectOnLine(row.in1,row.in2)
			require.Equal(t, row.out, res.RoundPlaces(8))
		})
	}
}
