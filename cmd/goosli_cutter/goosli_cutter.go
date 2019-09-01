package main

import (
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

var (
	stl  = kingpin.Flag("stl", "Source stl file to cut.").Short('s').Required().String()
	out1 = kingpin.Flag("out1", "Output path for upper result.").Short('o').Default("out1.stl").String()
	out2 = kingpin.Flag("out2", "Output path for down result.").Short('u').Default("out2.stl").String()
	px   = kingpin.Flag("pointx", "Plane's point x coord.").Short('x').Default("0").Float64()
	py   = kingpin.Flag("pointy", "Plane's point y coord.").Short('y').Default("0").Float64()
	pz   = kingpin.Flag("pointz", "Plane's point z coord.").Short('z').Default("0").Float64()
	nx   = kingpin.Flag("normali", "Plane's normal x coord.").Short('i').Default("0").Float64()
	ny   = kingpin.Flag("normalj", "Plane's normal y coord.").Short('j').Default("0").Float64()
	nz   = kingpin.Flag("normalk", "Plane's normal z coord.").Short('k').Default("1").Float64()
)

func main() {

	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}

	mUp, mDown, err := helpers.CutMesh(mesh, Plane{Point{*px, *py, *pz}, V(*nx, *ny, *nz)})
	if err != nil {
		log.Fatal("failed to cut mesh: ", err)
	}

	err = SaveSTL(*out1, mUp)
	if err != nil {
		log.Fatal("failed to save upper mesh: ", err)
	}

	err = SaveSTL(*out2, mDown)
	if err != nil {
		log.Fatal("failed to save upper mesh: ", err)
	}
}
