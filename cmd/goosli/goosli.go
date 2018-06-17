package main

import (
	"log"
	"github.com/l1va/goosli/slicers"
	"io/ioutil"
	. "github.com/l1va/goosli/primitives"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
)

var (
	stl       = kingpin.Flag("stl", "Source stl file to slice.").Short('s').Required().String()
	gcode     = kingpin.Flag("gcode", "Output file for gcode.").Short('o').Default("out.gcode").String()
	ox        = kingpin.Flag("originx", "Center of model by x.").Short('x').Default("0").Float64()
	oy        = kingpin.Flag("originy", "Center of model by y.").Short('y').Default("0").Float64()
	oz        = kingpin.Flag("originz", "Min of model by z.").Short('z').Default("0").Float64()

	epsilon   = kingpin.Flag("epsilon", "Simplification line parameter.").Short('e').Default("10.0").Float64()

	thickness = kingpin.Flag("thickness", "Set the slice thickness.").Short('t').Default("0.2").Float64()
	wallThickness = kingpin.Flag("wall_thickness", "Set the wall thickness.").Default("0.8").Float64()
	fillDensity = kingpin.Flag("fill_density", "Fill density in percents.").Default("20").Int()
	bedTemperature = kingpin.Flag("bed_temperature", "Bed temperature in Celsius.").Default("60").Int()
	extruderTemperature = kingpin.Flag("extruder_temperature", "Extruder temperature in Celsius.").Default("200").Int()
	travelSpeed = kingpin.Flag("travel_speed", "Printing travel speed.").Default("9000").Int()
)

func settings() slicers.Settings {
	return slicers.Settings{
	DateTime :time.Now().Format(time.RFC822),
	LayerHeight :*thickness,
	WallThickness:*wallThickness,
	FillDensity : *fillDensity,
	BedTemperature: *bedTemperature,
	ExtruderTemperature :*extruderTemperature,
	TravelSpeed : *travelSpeed,
	LayerCount : 0,
	}
}

func main() {

	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}
	mesh.Shift(V(-*ox, -*oy, -*oz))

	buffer := slicers.SliceByProfile(mesh, *epsilon, settings())

	err = ioutil.WriteFile(*gcode, buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save gcode in file: ", err)
	}
}
