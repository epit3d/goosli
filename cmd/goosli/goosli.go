package main

import (
	"log"
	"github.com/l1va/goosli/slicers"
	"io/ioutil"
	. "github.com/l1va/goosli/primitives"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
	"bytes"
	"github.com/l1va/goosli/slicers/vip"
	"github.com/l1va/goosli/gcode"
)

var (
	stl = kingpin.Flag("stl", "Source stl file to slice.").Short('s').Required().String()
	out = kingpin.Flag("gcode", "Output file for gcode.").Short('o').Default("out.gcode").String()
	ox  = kingpin.Flag("originx", "Center of model by x.").Short('x').Default("0").Float64()
	oy  = kingpin.Flag("originy", "Center of model by y.").Short('y').Default("0").Float64()
	oz  = kingpin.Flag("originz", "Min of model by z.").Short('z').Default("0").Float64()

	epsilon = kingpin.Flag("epsilon", "Simplification line parameter.").Short('e').Default("10.0").Float64()

	thickness           = kingpin.Flag("thickness", "Set the slice thickness.").Short('t').Default("0.2").Float64()
	wallThickness       = kingpin.Flag("wall_thickness", "Set the wall thickness.").Default("1.2").Float64()
	fillDensity         = kingpin.Flag("fill_density", "Fill density in percents.").Default("20").Int()
	bedTemperature      = kingpin.Flag("bed_temperature", "Bed temperature in Celsius.").Default("60").Int()
	extruderTemperature = kingpin.Flag("extruder_temperature", "Extruder temperature in Celsius.").Default("200").Int()
	printSpeed          = kingpin.Flag("print_speed", "Printing speed.").Default("50").Int()
	nozzle              = kingpin.Flag("nozzle", "Nozzle diameter.").Default("0.4").Float64()

	slicingType = kingpin.Flag("slicing_type", "Slicing type.").Default("vip").String()
)

func settings() slicers.Settings {
	return slicers.Settings{
		DateTime:            time.Now().Format(time.RFC822),
		Epsilon:             *epsilon,
		LayerHeight:         *thickness,
		WallThickness:       *wallThickness,
		FillDensity:         *fillDensity,
		BedTemperature:      *bedTemperature,
		ExtruderTemperature: *extruderTemperature,
		PrintSpeed:          *printSpeed * 60,
		Nozzle:              *nozzle,
		LayerCount:          0,
	}
}

func main() {

	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}
	//most := V(-60.08446554467082, -35.0, 0.0)
	mesh.Shift(V(*ox, *oy, *oz))
	//mesh.Shift(most)

	var gcd gcode.Gcode
	setts := settings()

	if *slicingType == "3axes" {
		gcd = slicers.SliceByVectorToGcode(mesh, AxisZ, setts)
	} else if *slicingType == "5axes_by_profile" {
		gcd = slicers.SliceByProfile(mesh, setts)
	} else if *slicingType == "5axes" {
		gcd = slicers.Slice5Axes(mesh, setts)
	} else if *slicingType == "vip" {
		gcd = vip.Slice(mesh, setts)
	} else {
		log.Fatal("unsupported slicing type: ", *slicingType)
	}

	buf := CommandsWithTemplates(gcd, setts)

	err = ioutil.WriteFile(*out, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save gcode in file: ", err)
	}
}

func CommandsWithTemplates(gcd gcode.Gcode, settings slicers.Settings) bytes.Buffer {
	settings.LayerCount = gcd.LayersCount
	smap := settings.ToMap()

	var buffer bytes.Buffer
	buffer.WriteString(PrepareDataFile("data/header_template.txt", smap))
	for _, cmd := range gcd.Cmds {
		cmd.ToGCode(&buffer)
	}
	buffer.WriteString(PrepareDataFile("data/footer_template.txt", smap))
	return buffer
}
