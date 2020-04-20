package slicers

import (
	"github.com/l1va/goosli/gcode"
	"strconv"

	. "github.com/l1va/goosli/primitives"
)

type Settings struct {
	GcodeSettings *gcode.GcodeSettings
	DateTime            string
	Epsilon             float64
	LayerHeight         float64
	WallThickness       float64
	FillDensity         int
	BedTemperature      int
	ExtruderTemperature int
	LayerCount          int
	RotationCenterZ     float64
	PlanesFile          string
	FillingType         string
	ColorizedAngle      float64
	UnitVector          Vector
	SupportsOn          bool
	SupportOffset       float64
}

func (s *Settings) ToMap() map[string]string {
	return map[string]string{
		"{datetime}":             s.DateTime,
		"{layer_height}":         StrF(s.LayerHeight),
		"{wall_thickness}":       StrF(s.WallThickness),
		"{fill_density}":         strconv.Itoa(s.FillDensity),
		"{bed_temperature}":      strconv.Itoa(s.BedTemperature),
		"{extruder_temperature}": strconv.Itoa(s.ExtruderTemperature),
		"{print_speed}":          strconv.Itoa(s.GcodeSettings.PrintSpeed),
		//"{nozzle}":               StrF(s.LineWidth),
		"{layer_count}":          strconv.Itoa(s.LayerCount),
	}
}
