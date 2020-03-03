package vip

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
)

type AnalyzedPlane struct {
	Plane
	tilted bool
	rotz   float64
}

func (p AnalyzedPlane) String() string {
	return fmt.Sprintf("%s T%s R%s", p.P.String(), strconv.FormatBool(p.tilted), StrF(p.rotz))
}

func ParsePlane(line string) AnalyzedPlane {
	p := AnalyzedPlane{}
	words := strings.Fields(line)
	var err error
	if p.tilted, err = strconv.ParseBool(words[3][1:]); err != nil {
		fmt.Printf("ERROR parsing bool: %v\n", err)
	}
	p.rotz = ParseFloat(words[4][1:])
	v := AxisZ
	if p.tilted {
		v = v.RotateAbout(AxisX, -angleX)
	}
	v = v.RotateAbout(AxisZ, p.rotz)
	p.Plane = Plane{
		P: Point{
			X: ParseFloat(words[0][1:]),
			Y: ParseFloat(words[1][1:]),
			Z: ParseFloat(words[2][1:]),
		}, N: v,
	}
	return p
}

// Analyze - Analyzing and suggesting plane rotations for future slicing
func Analyze(mesh *Mesh, angle float64) []AnalyzedPlane {
	ans := []AnalyzedPlane{}

	curPlane := AnalyzedPlane{PlaneXY, false, 0}
	curPlaneV := curPlane.N
	curMesh := mesh

	for {
		if len(ans) > 5 || len(curMesh.Triangles) == 0 { //TODO: is it good way to remove possibility of infinite loop?
			return ans
		}

		fails := helpers.ColorizeTriangles(*curMesh, angle, curPlaneV)
		var t *Triangle
		for i, fail := range fails {
			if fail {
				t2 := mesh.Triangles[i]
				if t2.P1.Z < 0.01 && t2.P2.Z < 0.01 && t2.P3.Z < 0.01 { // it is first layer, skip
					continue
				}
				if t == nil || t.MinZ(curPlaneV) > t2.MinZ(curPlaneV) {
					t = &t2
				}
			}
		}
		if t == nil {
			return ans
		}
		minpr := t.P1.ToVector().Dot(curPlaneV)
		pr2 := t.P2.ToVector().Dot(curPlaneV)
		pr3 := t.P3.ToVector().Dot(curPlaneV)
		p := t.P1
		if pr2 < minpr {
			p = t.P2
			minpr = pr2

		}
		if pr3 < minpr {
			p = t.P3
			minpr = pr3
		}

		if curPlane.tilted {
			curPlane = AnalyzedPlane{Plane{P: p, N: AxisZ}, false, 0}
		} else {
			angle := t.N.ProjectOnPlane(PlaneXY).Angle(AxisX) - 90
			a := t.N.ProjectOnPlane(PlaneXY)
			b := a.Rotate(RotationAroundZ(90))
			norm := a.RotateAbout(b, 90-angleX)
			println("angle x:", angle, norm.Angle(AxisX), norm.ProjectOnPlane(PlaneXY).Angle(AxisX))
			curPlane = AnalyzedPlane{Plane{P: p, N: norm}, true, angle}
		}
		ans = append(ans, curPlane)
		curPlaneV = curPlane.N
		var err error
		curMesh, _, err = helpers.CutMesh(curMesh, curPlane.Plane)
		if err != nil {
			log.Fatal("failed to cut mesh, during analization: ", err)
		}
		//if curPlane.tilted {

		//}else{
		//	println("POINT to CUT: ",curPlane.P.String())
		//SaveSTL("cuttedSTL.stl", curMesh)
		//return ans
	}
}
