package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"time"

	"github.com/l1va/goosli/gcode"
	. "github.com/l1va/goosli/primitives"
	"github.com/l1va/goosli/slicers"
	"github.com/l1va/goosli/slicers/vip"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	stl = kingpin.Flag("stl", "Source stl file to slice.").Short('s').Required().String()
	out = kingpin.Flag("gcode", "Output file for gcode.").Short('o').Default("out.gcode").String()
	ox  = kingpin.Flag("originx", "Center of model by x.").Short('x').Default("0").Float64()
	oy  = kingpin.Flag("originy", "Center of model by y.").Short('y').Default("0").Float64()
	oz  = kingpin.Flag("originz", "Min of model by z.").Short('z').Default("0").Float64()
	rcx = kingpin.Flag("rotcx", "X coord of rotation center.").Default("0").Float64()
	rcy = kingpin.Flag("rotcy", "Y coord of rotation center.").Default("0").Float64()
	rcz = kingpin.Flag("rotcz", "Z coord of rotation center.").Default("0").Float64()

	epsilon = kingpin.Flag("epsilon", "Simplification line parameter.").Short('e').Default("10.0").Float64()

	layerHeight         = kingpin.Flag("layer_height", "Set the slice layer height.").Short('t').Default("0.2").Float64()
	wallThickness       = kingpin.Flag("wall_thickness", "Set the wall thickness.").Default("1.2").Float64()
	fillDensity         = kingpin.Flag("fill_density", "Fill density in percents.").Default("20").Int()
	topLayers           = kingpin.Flag("top_layers", "Number of top layers").Default("3").Int()
	bottomLayers        = kingpin.Flag("bottom_layers", "Number of bottom layers").Default("3").Int()
	bedTemperature      = kingpin.Flag("bed_temperature", "Bed temperature in Celsius.").Default("60").Int()
	extruderTemperature = kingpin.Flag("extruder_temperature", "Extruder temperature in Celsius.").Default("200").Int()
	printSpeed          = kingpin.Flag("print_speed", "Printing speed.").Default("50").Int()
	printSpeedLayer1    = kingpin.Flag("print_speed_layer1", "Printing speed of Layer 1.").Default("50").Int()
	printSpeedWall      = kingpin.Flag("print_speed_wall", "Printing speed of walls.").Default("50").Int()
	fanOffLayer1        = kingpin.Flag("fan_off_layer1", "Turn off the fan for Layer 1.").Bool()
	lineWidth           = kingpin.Flag("line_width", "Line Width.").Default("0.4").Float64()
	fillingType         = kingpin.Flag("filling_type", "Filling type(Lines,Squares,Triangles)").Default("Lines").String()
	retraction          = kingpin.Flag("retraction_on", "Turn on the retraction.").Bool()
	retractionSpeed     = kingpin.Flag("retraction_speed", "How fast to pull in the fillament.").Int()
	retractionDistance  = kingpin.Flag("retraction_distance", "How much fillament to pull in.").Float64()

	planesFile    = kingpin.Flag("planes_file", "File with planes description.").Default("planes_file.txt").String()
	slicingType   = kingpin.Flag("slicing_type", "Slicing type.").Default("vip").String()
	angle         = kingpin.Flag("angle", "Angle of bias to colorize a triangle.").Short('a').Default("30").Float64()
	nx            = kingpin.Flag("normali", "Plane's normal x coord.").Short('i').Default("0").Float64()
	ny            = kingpin.Flag("normalj", "Plane's normal y coord.").Short('j').Default("0").Float64()
	nz            = kingpin.Flag("normalk", "Plane's normal z coord.").Short('k').Default("1").Float64()
	supportsOn    = kingpin.Flag("supports_on", "Add supports").Bool()
	supportOffset = kingpin.Flag("support_offset", "Offset (shifting) for support").Default("1.0").Float64()

	barDiameter = kingpin.Flag("bar_diameter", "Plastic bar diameter").Default("1.75").Float64()
	flow        = kingpin.Flag("flow", "Printing flow (0;1]").Default("1.0").Float64()

	skirtLineCount = kingpin.Flag("skirt_line_count", "Build plate Adhesion: skirt line count").Default("3").Int()
)

//TODO: create one binary, not 4

func settings() slicers.Settings {
	gcodeSettings := gcode.GcodeSettings{
		BarDiameter:        *barDiameter,
		Flow:               *flow,
		LayerHeight:        *layerHeight,
		LineWidth:          *lineWidth,
		FanOffLayer1:       *fanOffLayer1,
		PrintSpeed:         *printSpeed * 60,
		PrintSpeedLayer1:   *printSpeedLayer1 * 60,
		PrintSpeedWall:     *printSpeedWall * 60,
		Retraction:         *retraction,
		RetractionSpeed:    *retractionSpeed,
		RetractionDistance: *retractionDistance,
	}

	return slicers.Settings{
		GcodeSettings:       &gcodeSettings,
		DateTime:            time.Now().Format(time.RFC822),
		Epsilon:             *epsilon,
		LayerHeight:         *layerHeight,
		WallThickness:       *wallThickness,
		FillDensity:         *fillDensity,
		TopLayers:           *topLayers,
		BottomLayers:        *bottomLayers,
		BedTemperature:      *bedTemperature,
		ExtruderTemperature: *extruderTemperature,
		LayerCount:          0,
		RotationCenter:      Point{*rcx, *rcy, *rcz},
		PlanesFile:          *planesFile,
		FillingType:         *fillingType,
		ColorizedAngle:      *angle,
		UnitVector:          V(*nx, *ny, *nz), //TODO: seems useless, recheck
		SupportsOn:          *supportsOn,
		SupportOffset:       *supportOffset,
		SkirtLineCount: *skirtLineCount,
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
	settings.LayerCount = gcd.LayerCount()
	smap := settings.ToMap()

	var buffer bytes.Buffer
	buffer.WriteString(PrepareDataFile("data/header_template.txt", smap))
	gcd.ToOutput(&buffer)
	buffer.WriteString(PrepareDataFile("data/footer_template.txt", smap))
	return buffer
}
