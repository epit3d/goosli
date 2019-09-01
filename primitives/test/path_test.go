package test

import (
	//"fmt"
	"testing"
	"fmt"
	//"github.com/stretchr/testify/require"
	. "github.com/l1va/goosli/primitives"
	"github.com/stretchr/testify/require"
)

func TestPath_Join(t *testing.T) {
	//cases := []struct {
	//	in  []Path
	//	out []Path
	//}{
	//	{
	//		in: []Path{{[]Line{{Point{1, 1, 1}, Point{1, 1, 2}}}},
	//			{[]Line{{Point{1, 1, 2}, Point{1, 1, 3}}}},
	//			{[]Line{{Point{1, 1, 3}, Point{1, 1, 4}}}},
	//			{[]Line{{Point{1, 1, 4}, Point{1, 1, 1}}}}},
	//		out: []Path{{[]Line{{Point{1, 1, 1}, Point{1, 1, 2}},
	//			{Point{1, 1, 2}, Point{1, 1, 3}},
	//			{Point{1, 1, 3}, Point{1, 1, 4}},
	//			{Point{1, 1, 4}, Point{1, 1, 1}}}}},
	//	},
	//}
	//for i, row := range cases {
	//	t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
	//		res := JoinPaths(row.in)
	//		require.Equal(t, len(row.out), len(res) )
	//		used := []bool{}
	//		for
	//		for _, p := range res {
	//			found := false
	//			for _, p2 := range row.out {
	//				if p.Equal(p2) {
	//					found = true
	//				}
	//			}
	//			require.Equal(t, len(p), len(row.in) )
	//
	//		}
	//		require.Equal(t, row.out, )
	//	})
	//}
}

func TestPath_FindCentroid(t *testing.T) {
	cases := []struct {
		in  Path
		out Point
	}{
		{
			in: Path{Lines: []Line{{P1: Point{1, 1, 1}, P2: Point{1, 1, 2}},
				{P1: Point{1, 1, 2}, P2: Point{1, 2, 2}},
				{P1: Point{1, 2, 2}, P2: Point{1, 2, 1}},
				{P1: Point{X: 1, Y: 2, Z: 1}, P2: Point{X: 1, Y: 1, Z: 1}},}},
			out: Point{X: 1, Y: 1.5, Z: 1.5},
		},
		{
			in: Path{Lines: []Line{{Point{2, 1, 1}, Point{2, 1, 2}},
				{Point{2, 1, 2}, Point{1, 2, 2}},
				{Point{1, 2, 2}, Point{1, 2, 1}},
				{Point{1, 2, 1}, Point{2, 1, 1}},}},
			out: Point{1.5, 1.5, 1.5},
		},
		{
			in: Path{Lines: []Line{{Point{2, 1, 1}, Point{2, 2, 2}},
				{Point{2, 2, 2}, Point{1, 3, 1}},
				{Point{1, 3, 1}, Point{1, 2, 3}},
				{Point{1, 2, 3}, Point{2, 1, 1}},}},
			out: Point{1.5, 2, 2},
		},
		{
			in: Path{Lines: []Line{{Point{15, 3, 2}, Point{3, 2, 12}},
				{Point{3, 2, 12}, Point{1, 13, 10}},
				{Point{1, 13, 10}, Point{12, 11, 3}},
				{Point{12, 11, 3}, Point{15, 3, 2}},}},
			out: Point{X: 7.472222222222222, Y: 7.138888888888889, Z: 7.444444444444445},
		},
	}
	for i, row := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require.Equal(t, row.out, FindCentroid(row.in))
		})
	}
}
