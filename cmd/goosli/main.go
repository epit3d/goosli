package main

import (
	"log"
	"github.com/l1va/goosli/slicers"
	"bytes"
	"io/ioutil"
	. "github.com/l1va/goosli"
	"strconv"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	stl       = kingpin.Flag("stl", "Source stl file to slice.").Short('s').Required().String()
	gcode     = kingpin.Flag("gcode", "Output file for gcode.").Short('o').Default("out.gcode").String()
	thickness = kingpin.Flag("thickness", "Set the slice thickness.").Short('t').Default("0.2").Float64()
	alpha     = kingpin.Flag("alpha", "Angle of slicing rotation in degrees.").Short('a').Default("30").Float64()
)

func main() {
	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}
	cmds := slicers.SliceByZ(*mesh, *thickness, V(1,1,3))
	//cmds := slicers.SliceWithSlope(*mesh, *thickness, *alpha)

	var buffer bytes.Buffer
	for i := 0; i < len(cmds); i++ {
		buffer.WriteString(";Layer" + strconv.Itoa(i) + "\n")
		buffer.WriteString(cmds[i].ToGCode())
	}

	err = ioutil.WriteFile(*gcode, buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save gcode in file: ", err)
	}
}
