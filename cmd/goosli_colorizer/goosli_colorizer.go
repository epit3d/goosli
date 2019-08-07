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
	angle           = kingpin.Flag("angle", "Angle of bias to colorize a triangle.").Short('a').Default("30").Float64()
	out = kingpin.Flag("out", "Output path for result.").Short('o').Default("colorize_triangles.txt").String()
)

func main() {

	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}

	arr := helpers.ColorizeTriangles(*mesh, *angle)

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
