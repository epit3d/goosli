package gcode

import (
	"bytes"
	. "github.com/l1va/goosli/primitives"
)

type GcodeSettings struct {
	BarDiameter        float64
	Flow               float64
	LayerHeight        float64
	LineWidth          float64
	FanOffLayer1       bool
	PrintSpeed         int
	PrintSpeedLayer1   int
	PrintSpeedWall     int
	Retraction         bool
	RetractionSpeed    int
	RetractionDistance float64
}
type Gcode struct {
	Settings GcodeSettings
	Cmds     []Command
}

func NewGcode(settings GcodeSettings) Gcode {
	return Gcode{
		Settings: settings,
	}
}

func (g *Gcode) Add(cmd Command) {
	g.Cmds = append(g.Cmds, cmd)
}

func (g *Gcode) AddLayers(lays []Layer) {
	g.Cmds = append(g.Cmds, LayersMoving{Layers: lays})
}

func (g *Gcode) ToOutput(b *bytes.Buffer) {
	state := NewState()
	for _, cmd := range g.Cmds {
		state = cmd.ToGCode(b, state, g.Settings)
	}
}

func (g *Gcode) LayerCount() int {
	c := 0
	for _, cmd := range g.Cmds {
		c += cmd.LayersCount()
	}
	return c
}
