package main

import (
	"log"
	"github.com/l1va/goosli/slicers"
	"io/ioutil"
	"github.com/l1va/goosli"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	stl       = kingpin.Flag("stl", "Source stl file to slice.").Short('s').Required().String()
	gcode     = kingpin.Flag("gcode", "Output file for gcode.").Short('o').Default("out.gcode").String()
	thickness = kingpin.Flag("thickness", "Set the slice thickness.").Short('t').Default("0.2").Float64()
	epsilon   = kingpin.Flag("epsilon", "Simplification line parameter.").Short('e').Default("10.0").Float64()
	ox        = kingpin.Flag("originx", "Shift center of model by x.").Short('x').Default("0").Float64()
	oy        = kingpin.Flag("originy", "Shift center of model by y.").Short('y').Default("0").Float64()
	oz        = kingpin.Flag("originz", "Shift center of model by z.").Short('z').Default("0").Float64()
)

func main() {

	kingpin.Parse()

	mesh, err := goosli.LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}
	mesh.Shift(goosli.V(*ox, *oy, *oz))

	buffer := slicers.SliceByProfile(mesh, *thickness, *epsilon)

	err = ioutil.WriteFile(*gcode, buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save gcode in file: ", err)
	}
}
