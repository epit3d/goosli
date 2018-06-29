package test

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/l1va/goosli/slicers"
)

func TestSlicer_MiddleOnTheRing(t *testing.T) {
	cases := []struct {
		in  []float64
		in2 int
		out int
	}{
		{
			in:  []float64{10, 20, 30, 40},
			in2: 360,
			out: 25,
		},
		{
			in:  []float64{50, 60, 300, 310},
			in2: 360,
			out: 0,
		},
		{
			in:  []float64{50, 60, 70, 300, 310},
			in2: 360,
			out: 14,
		},
		{
			in:  []float64{10, 20, 170, 190, 340, 350}, // ???
			in2: 360,
			out: 300,
		},
	}

	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require.Equal(t, row.out, slicers.MiddleOnTheRing(row.in, row.in2))
		})
	}
}
