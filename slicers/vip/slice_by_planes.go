package vip

import (
	"github.com/l1va/goosli/debug"
	"github.com/l1va/goosli/gcode"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	"github.com/l1va/goosli/slicers"
	"log"
)

type BedPlane struct {
	tilted bool
	rotz   float64
}

func SliceByPlanes(mesh *Mesh, settings slicers.Settings, cutPlanes []AnalyzedPlane) gcode.Gcode {
	gcd := gcode.NewGcode(*settings.GcodeSettings)
	debug.RecreateFile()
	fillPlanes, fullFillPlanes := slicers.CalcFillPlanes(mesh, settings)
	var up *Mesh
	var err error

	parts := make([]*Mesh, len(cutPlanes)+1)
	for i := len(cutPlanes) - 1; i >= 0; i-- { //cut from the last plane (in reverse order)
		up, mesh, err = helpers.CutMesh(mesh, cutPlanes[i].Plane)
		if err != nil {
			log.Fatal("failed to cut mesh, by plane: ", err, cutPlanes[i].Plane)
		}
		parts[i+1] = up
	}
	parts[0] = mesh //do not forget first part

	bedPlanes := calcBedPlanes(cutPlanes) // len(bedPlanes) = cutplanes + 1
	rotated := false
	for i := 0; i < len(bedPlanes); i++ {
		pl := bedPlanes[i]
		curmesh := parts[i]

		if i != 0 {
			curmesh = curmesh.Rotate(RotationAroundZ(pl.rotz), settings.RotationCenter)
			gcd.Add(gcode.RotateZ{Angle: pl.rotz})
		}
		if pl.tilted {
			curmesh = curmesh.Rotate(RotationAroundX(-angleX), settings.RotationCenter)
			gcd.Add(gcode.InclineX{})
			rotated = true
		}
		if rotated && !pl.tilted {
			gcd.Add(gcode.InclineXBack{})
			rotated = false
		}
		add := slicers.SliceByVector(curmesh, AxisZ, settings)
		filled := PrepareLayers(add, settings, fillPlanes, fullFillPlanes) //TODO: fillPlanes fix
		if i == 0 && len(filled) > 0 {
			filled[0] = SkirtPathes(filled[0], settings.SkirtLineCount, settings.GcodeSettings.LineWidth)
		}
		gcd.AddLayers(filled)
	}

	return gcd
}

func calcBedPlanes(cutPlanes []AnalyzedPlane) []BedPlane {
	res := []BedPlane{BedPlane{tilted: false, rotz: 0}}
	for _, p := range cutPlanes {
		res = append(res, BedPlane{tilted: p.tilted, rotz: p.rotz})
	}
	return res
}
