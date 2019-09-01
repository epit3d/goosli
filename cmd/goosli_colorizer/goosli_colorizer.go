package main

import (
	"bufio"
	"encoding/binary"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

var (
	stl = kingpin.Flag("stl", "Source stl file to colorize.").Short('s').Required().String()
	//stl = kingpin.Flag("stl", "Source stl file to colorize.").Short('s').Default("/home/l1va/Downloads/bunny.stl").String()
	angle = kingpin.Flag("angle", "Angle of bias to colorize a triangle.").Short('a').Default("30").Float64()
	out   = kingpin.Flag("out", "Output path for result.").Short('o').Default("colorize_triangles.txt").String()
	nx    = kingpin.Flag("normali", "Plane's normal x coord.").Short('i').Default("0").Float64()
	ny    = kingpin.Flag("normalj", "Plane's normal y coord.").Short('j').Default("0").Float64()
	nz    = kingpin.Flag("normalk", "Plane's normal z coord.").Short('k').Default("1").Float64()
)

func main() {

	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}

	vn := V(*nx, *ny, *nz)
	arr := helpers.ColorizeTriangles(*mesh, *angle, vn)

	file, err := os.Create(*out)
	if err != nil {
		log.Fatal("failed to save colorize triangles: ", err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)

	err = binary.Write(w, binary.LittleEndian, arr)
	if err != nil {
		log.Fatal("failed to write value of colorize triangles: ", err)
	}

	w.Flush()
}
