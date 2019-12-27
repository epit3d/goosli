package main

import (
	"bytes"
	"fmt"
	. "github.com/l1va/goosli/primitives"
	"github.com/l1va/goosli/slicers/vip"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
)

var (
	stl = kingpin.Flag("stl", "Source stl file to colorize.").Short('s').Required().String()
	//stl = kingpin.Flag("stl", "Source stl file to colorize.").Short('s').Default("/home/l1va/Downloads/bunny.stl").String()
	out = kingpin.Flag("out", "Output path for result.").Short('o').Default("analyzed_planes.txt").String()
	angle = kingpin.Flag("angle", "Angle of bias to colorize a triangle.").Short('a').Default("30").Float64()
	ox  = kingpin.Flag("originx", "Center of model by x.").Short('x').Default("0").Float64()
	oy  = kingpin.Flag("originy", "Center of model by y.").Short('y').Default("0").Float64()
	oz  = kingpin.Flag("originz", "Min of model by z.").Short('z').Default("0").Float64()
	rcx = kingpin.Flag("rotcx", "X coord of rotation center.").Default("0").Float64()
	rcy = kingpin.Flag("rotcy", "Y coord of rotation center.").Default("0").Float64()
	rcz = kingpin.Flag("rotcz", "Z coord of rotation center.").Default("0").Float64()
)

func main() {

	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}
	mesh.Shift(V(*ox, *oy, *oz))

	planes := vip.Analyze(mesh, *angle)

	var buffer bytes.Buffer
	println("planes:",len(planes), planes)
	for _, pl := range planes {
		buffer.WriteString(fmt.Sprintf("%s\n", pl))
	}

	err = ioutil.WriteFile(*out, buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save gcode in file: ", err)
	}

}
