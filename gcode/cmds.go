package gcode

import (
	"bytes"
	"math"
	"strconv"

	. "github.com/l1va/goosli/primitives"
)

type State struct {
	fan        bool
	layerIndex int
	eOff       float64
	feedrate   int
	isWall     bool
	pos        Point
}

func NewState() State {
	return State{
		fan:        false, //start from fan off
		layerIndex: 0,
		eOff:       0,
		feedrate:   0,
		isWall:     false,
		pos:        OriginPoint,
	}
}

type Command interface {
	ToGCode(b *bytes.Buffer, state State, sett GcodeSettings) State
	LayersCount() int
}

type InclineXBack struct {
}

func (r InclineXBack) ToGCode(b *bytes.Buffer, state State, sett GcodeSettings) State {
	b.WriteString("G0 X0 Y0 Z200 \n")
	b.WriteString("T2 \n")
	b.WriteString("G92 E0 \n")
	b.WriteString("G1 F300 E0 \n")
	b.WriteString("G92 E0 \n")
	b.WriteString("T0 \n")
	b.WriteString("G92 E" + StrF(state.eOff) + "\n")
	return state
}
func (r InclineXBack) LayersCount() int {
	return 0
}

type InclineX struct {
}

func (r InclineX) ToGCode(b *bytes.Buffer, state State, sett GcodeSettings) State {
	b.WriteString("G0 X0 Y0 Z200 \n")
	b.WriteString("T2 \n")
	b.WriteString("G92 E0 \n")
	b.WriteString("G1 F300 E60 \n")
	b.WriteString("G92 E0 \n")
	b.WriteString("T0 \n")
	b.WriteString("G92 E" + StrF(state.eOff) + "\n")
	return state

}
func (r InclineX) LayersCount() int {
	return 0
}

type RotateXZ struct {
	AngleX float64
	AngleZ float64
}

func (r RotateXZ) ToGCode(b *bytes.Buffer, st State, sett GcodeSettings) State {
	b.WriteString("G62 X" + StrF(r.AngleX) + " Z" + StrF(r.AngleZ) + "\n")
	return st

}
func (r RotateXZ) LayersCount() int {
	return 0
}

type RotateZ struct {
	Angle float64
}

func (r RotateZ) ToGCode(b *bytes.Buffer, state State, sett GcodeSettings) State {
	//b.WriteString("G0 A" + StrF(r.Angle) + "\n")
	b.WriteString("G0 X0 Y0 Z200 \n")
	b.WriteString("T1 \n")
	b.WriteString("G92 E0 \n")
	b.WriteString("G1 F300 E" + StrF(r.Angle) + " \n")
	b.WriteString("G92 E0 \n")
	b.WriteString("T0 \n")
	b.WriteString("G92 E" + StrF(state.eOff) + "\n")
	return state

}
func (r RotateZ) LayersCount() int {
	return 0
}

type LayersMoving struct {
	Layers []Layer
}

func (lm LayersMoving) ToGCode(b *bytes.Buffer, state State, sett GcodeSettings) State {
	for i := 0; i < len(lm.Layers); i++ {
		b.WriteString(";LAYER:" + strconv.Itoa(state.layerIndex) + "\n")
		state = switchFanGCode(state, sett, b)

		state.isWall = true
		state = pathesToGCode(lm.Layers[i].Paths, "OUTER_PATHES", sett, state, b)
		state = pathesToGCode(lm.Layers[i].MiddlePs, "MIDDLE_PATHES", sett, state, b)
		state = pathesToGCode(lm.Layers[i].InnerPs, "INNER_PATHES", sett, state, b)
		state.isWall = false
		state = pathesToGCode(lm.Layers[i].Fill, "FILL_PATHES", sett, state, b)

		state = decreaseEOff(state, b)
		state.layerIndex += 1
	}
	return state
}

func (lm LayersMoving) LayersCount() int {
	return len(lm.Layers)
}

func decreaseEOff(state State, b *bytes.Buffer) State {
	if state.eOff > 4000 {
		b.WriteString("G92 E0\n")
		state.eOff = 0
	}
	return state
}

func switchFanGCode(state State, sett GcodeSettings, b *bytes.Buffer) State {
	if state.layerIndex == 0 && sett.FanOffLayer1 {
		if state.fan {
			b.WriteString("M107\n") //fan off
			state.fan = false
		}
	} else {
		if !state.fan {
			b.WriteString("M106\n")
			state.fan = true
		}
	}
	return state
}

func printSpeedToGCode(state State, sett GcodeSettings, b *bytes.Buffer) State {
	needed := state.feedrate
	if state.layerIndex == 0 {
		needed = sett.PrintSpeedLayer1
	} else {
		if state.isWall {
			needed = sett.PrintSpeedWall
		} else {
			needed = sett.PrintSpeed
		}
	}
	if needed != state.feedrate {
		b.WriteString("G0 F" + strconv.Itoa(needed) + "\n")
		state.feedrate = needed
	}
	return state
}

func retractionToGCode(b *bytes.Buffer, sett GcodeSettings) {
	if !sett.Retraction {
		return
	}

	b.WriteString("G1 F" + strconv.Itoa(sett.RetractionSpeed) + " E" + StrF(-sett.RetractionDistance) + "\n")
}

func pathesToGCode(pths []Path, comment string, sett GcodeSettings, state State, b *bytes.Buffer) State {
	b.WriteString(";" + comment + "\n")

	// Set the printing speed for this path
	state = printSpeedToGCode(state, sett, b)

	for _, p := range pths {
		// Retraction first
		b.WriteString("G0 " + state.pos.StringDeltaZ(0.5) + "\n") //z up a little during retract
		retractionToGCode(b, sett)
		b.WriteString("G0 " + p.Points[0].StringDelta(0, 0, 0.5) + "\n")
		b.WriteString("G0 " + p.Points[0].StringDeltaZ(0) + "\n")

		for i := 1; i < len(p.Points); i++ {
			p1 := p.Points[i-1]
			p2 := p.Points[i]
			lDist := math.Sqrt(math.Pow(p2.X-p1.X, 2) + math.Pow(p2.Y-p1.Y, 2) + math.Pow(p2.Z-p1.Z, 2))
			state.eOff += (4 * sett.LineWidth * sett.LayerHeight * lDist) / (math.Pow(sett.BarDiameter, 2) * math.Pi)
			b.WriteString("G1 " + p2.StringGcode(p1) + " E" + StrF(state.eOff) + "\n")
		}
		state.pos = p.Points[len(p.Points)-1]
	}
	return state
}
