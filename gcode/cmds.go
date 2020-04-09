package gcode

import (
	"bytes"
	"math"
	"strconv"

	. "github.com/l1va/goosli/primitives"
)

type Command interface {
	ToGCode(b *bytes.Buffer)
	LayersCount() int
}

type InclineXBack struct {
}

func (r InclineXBack) ToGCode(b *bytes.Buffer) {
	b.WriteString("M42 \n")
}
func (r InclineXBack) LayersCount() int {
	return 0
}

type InclineX struct {
}

func (r InclineX) ToGCode(b *bytes.Buffer) {
	b.WriteString("M43 \n")
}
func (r InclineX) LayersCount() int {
	return 0
}

type RotateXZ struct {
	AngleX float64
	AngleZ float64
}

func (r RotateXZ) ToGCode(b *bytes.Buffer) {
	b.WriteString("G62 X" + StrF(r.AngleX) + " Z" + StrF(r.AngleZ) + "\n")
}
func (r RotateXZ) LayersCount() int {
	return 0
}

type RotateZ struct {
	Angle float64
}

func (r RotateZ) ToGCode(b *bytes.Buffer) {
	b.WriteString("G0 A" + StrF(r.Angle) + "\n")
}
func (r RotateZ) LayersCount() int {
	return 0
}

type LayersMoving struct {
	Layers    []Layer
	Index     int
	ExtParams ExtrusionParams
}

func (lm LayersMoving) ToGCode(b *bytes.Buffer) {
	eOff := 0.0
	//reset extruder encoder
	b.WriteString("G92 E0\n")
	for i := 0; i < len(lm.Layers); i++ {
		b.WriteString(";LAYER:" + strconv.Itoa(i+lm.Index) + "\n")
		switchFanGCode(lm.Layers[i].FanOff, b)
		eOff = pathesToGCode(lm.Layers[i].Paths, "OUTER_PATHES", lm.Layers[i].WallPrintSpeed, lm.ExtParams, eOff, b)
		eOff = pathesToGCode(lm.Layers[i].MiddlePs, "MIDDLE_PATHES", lm.Layers[i].PrintSpeed, lm.ExtParams, eOff, b)
		eOff = pathesToGCode(lm.Layers[i].InnerPs, "INNER_PATHES", lm.Layers[i].PrintSpeed, lm.ExtParams, eOff, b)
		eOff = pathesToGCode(lm.Layers[i].Fill, "FILL_PATHES", lm.Layers[i].PrintSpeed, lm.ExtParams, eOff, b)
		eOff = decreaseEOff(eOff, b)
	}
}

func (lm LayersMoving) LayersCount() int {
	return len(lm.Layers)
}

func decreaseEOff(eOff float64, b *bytes.Buffer) float64 {
	if eOff > 4000 {
		b.WriteString("G92 E0\n")
		return 0.0
	} else {
		return eOff
	}
}

func switchFanGCode(fanOff bool, b *bytes.Buffer) {
	if fanOff {
		b.WriteString("M107\n")
	} else {
		b.WriteString("M106\n")
	}
}

func printSpeedToGCode(feedrate int, b *bytes.Buffer) {
	b.WriteString("G0 F" + strconv.Itoa(feedrate) + "\n")
}

func retractionToGCode(b *bytes.Buffer, retraction bool, retractionDistance float64, retractionSpeed int) {
	if !retraction {
		return
	}

	b.WriteString("; Retraction\n")
	b.WriteString("G1 F" + strconv.Itoa(retractionSpeed) + " E" + StrF(-retractionDistance) + "\n")
}

func pathesToGCode(pths []Path, comment string, feedrate int, extParams ExtrusionParams, eOff float64, b *bytes.Buffer) float64 {

	b.WriteString(";" + comment + "\n")

	// Set the printing speed for this path
	printSpeedToGCode(feedrate, b)

	for _, p := range pths {
		// Retraction first
		retractionToGCode(b, p.Retraction, p.RetractionDistance, p.RetractionSpeed)

		b.WriteString("G0 " + p.Points[0].String() + "\n")
		for i := 1; i < len(p.Points); i++ {
			p1 := p.Points[i-1]
			p2 := p.Points[i]
			lDist := math.Sqrt(math.Pow(p2.X-p1.X, 2) + math.Pow(p2.Y-p1.Y, 2) + math.Pow(p2.Z-p1.Z, 2))
			eOff += (4 * extParams.LineWidth * extParams.LayerHeight * lDist) / (math.Pow(extParams.BarDiameter, 2) * math.Pi)
			b.WriteString("G1 " + p2.StringGcode(p1) + " E" + StrF(eOff) + "\n")
		}
	}
	return eOff
}
