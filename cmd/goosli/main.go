package main

import (
	"log"
	"github.com/l1va/goosli/slicers"
	"bytes"
	"io/ioutil"
	. "github.com/l1va/goosli"
	"strconv"
	"gopkg.in/alecthomas/kingpin.v2"
	"math"
)

var (
	stl       = kingpin.Flag("stl", "Source stl file to slice.").Short('s').Required().String()
	gcode     = kingpin.Flag("gcode", "Output file for gcode.").Short('o').Default("out.gcode").String()
	thickness = kingpin.Flag("thickness", "Set the slice thickness.").Short('t').Default("0.2").Float64()
	ox        = kingpin.Flag("originx", "Shift center of model by x.").Short('x').Default("0").Float64()
	oy        = kingpin.Flag("originy", "Shift center of model by y.").Short('y').Default("0").Float64()
	oz        = kingpin.Flag("originz", "Shift center of model by z.").Short('z').Default("0").Float64()
)

func main() {

	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}
	buffer := slicers.Slice5a(mesh, *thickness)

	err = ioutil.WriteFile(*gcode, buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save gcode in file: ", err)
	}

}

func main3() {
	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}

	m2, m1, err := slicers.Cut(mesh, Plane{Point{0, 0, 30}, V(0, 0, 1)})
	if err != nil {
		log.Fatal("failed to make first cut mesh: ", err)
	}

	cmds := slicers.SliceByZ(*m1, *thickness, V(0, 0, 1))

	var buffer bytes.Buffer
	for i := 0; i < len(cmds); i++ {
		buffer.WriteString(";Layer" + strconv.Itoa(i) + "\n")
		buffer.WriteString(cmds[i].ToGCode())
	}

	alphaG := 30
	buffer.WriteString("G52 " + strconv.Itoa(alphaG) + "\n")

	c := Point{*ox, *oy, *oz}
	cv := c.ToVector()

	alpha := math.Pi * float64(alphaG) / 180.0
	// transposed matrix to rotate around X
	//mx := V(1, 0, 0)
	//my := V(0, math.Cos(alpha), math.Sin(alpha))
	//mz := V(0, -math.Sin(alpha), math.Cos(alpha))
	mx := V(math.Cos(alpha), math.Sin(alpha), 0)
	my := V(-math.Sin(alpha), math.Cos(alpha), 0)
	mz := V(0, 0, 1)

	triangles := make([]Triangle, len(m2.Triangles))
	rotatedMesh := NewMesh(triangles)
	for i, t := range m2.Triangles {
		p1 := c.VectorTo(t.P1).Rotate(mx, my, mz).Add(cv).ToPoint()
		p2 := c.VectorTo(t.P2).Rotate(mx, my, mz).Add(cv).ToPoint()
		p3 := c.VectorTo(t.P3).Rotate(mx, my, mz).Add(cv).ToPoint()
		rotatedMesh.Triangles[i].Fill(p1, p2, p3)
	}

	m3, m2, err := slicers.Cut(&rotatedMesh, Plane{Point{0, 0, 50}, V(0, 0, 1)})
	if err != nil {
		log.Fatal("failed to make second cut mesh: ", err)
	}

	cmds = slicers.SliceByZ(*m2, *thickness, V(0, 0, 1))

	for i := 0; i < len(cmds); i++ {
		buffer.WriteString(";Layer" + strconv.Itoa(i) + "\n")
		buffer.WriteString(cmds[i].ToGCode())
	}

	alphaG = 30
	buffer.WriteString("G42 " + strconv.Itoa(alphaG) + "\n")

	c = Point{*ox, *oy, *oz}
	cv = c.ToVector()

	alpha = math.Pi * float64(alphaG) / 180.0
	// transposed matrix to rotate around X
	mx = V(1, 0, 0)
	my = V(0, math.Cos(alpha), math.Sin(alpha))
	mz = V(0, -math.Sin(alpha), math.Cos(alpha))
	//mx = V(math.Cos(alpha), math.Sin(alpha), 0)
	//my = V(-math.Sin(alpha), math.Cos(alpha), 0 )
	//mz = V(0, 0, 1)

	triangles = make([]Triangle, len(m3.Triangles))
	rotatedMesh = NewMesh(triangles)
	for i, t := range m3.Triangles {
		p1 := c.VectorTo(t.P1).Rotate(mx, my, mz).Add(cv).ToPoint()
		p2 := c.VectorTo(t.P2).Rotate(mx, my, mz).Add(cv).ToPoint()
		p3 := c.VectorTo(t.P3).Rotate(mx, my, mz).Add(cv).ToPoint()
		rotatedMesh.Triangles[i].Fill(p1, p2, p3)
	}

	cmds = slicers.SliceByZ(rotatedMesh, *thickness, V(0, 0, 1))

	for i := 0; i < len(cmds); i++ {
		buffer.WriteString(";Layer" + strconv.Itoa(i) + "\n")
		buffer.WriteString(cmds[i].ToGCode())
	}

	err = ioutil.WriteFile(*gcode, buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save gcode in file: ", err)
	}

}
func main4() {
	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}

	m2, m1, err := slicers.Cut(mesh, Plane{Point{0, 0, 30}, V(0, 0, 1)})
	if err != nil {
		log.Fatal("failed to make first cut mesh: ", err)
	}

	cmds := slicers.SliceByZ(*m1, *thickness, V(1, 0, 1))

	var buffer bytes.Buffer
	for i := 0; i < len(cmds); i++ {
		buffer.WriteString(";Layer" + strconv.Itoa(i) + "\n")
		buffer.WriteString(cmds[i].ToGCode())
	}

	alphaG := 10
	buffer.WriteString("G42 " + strconv.Itoa(alphaG) + "\n")

	c := Point{*ox, *oy, *oz}
	cv := c.ToVector()

	alpha := math.Pi * float64(alphaG) / 180.0
	// transposed matrix to rotate around X
	mx := V(1, 0, 0)
	my := V(0, math.Cos(alpha), math.Sin(alpha))
	mz := V(0, -math.Sin(alpha), math.Cos(alpha))
	//mx := V(math.Cos(alpha), math.Sin(alpha), 0)
	//my := V(-math.Sin(alpha), math.Cos(alpha), 0 )
	//mz := V(0, 0, 1)

	triangles := make([]Triangle, len(m2.Triangles))
	rotatedMesh := NewMesh(triangles)
	for i, t := range m2.Triangles {
		p1 := c.VectorTo(t.P1).Rotate(mx, my, mz).Add(cv).ToPoint()
		p2 := c.VectorTo(t.P2).Rotate(mx, my, mz).Add(cv).ToPoint()
		p3 := c.VectorTo(t.P3).Rotate(mx, my, mz).Add(cv).ToPoint()
		rotatedMesh.Triangles[i].Fill(p1, p2, p3)
	}

	m3, m2, err := slicers.Cut(&rotatedMesh, Plane{Point{0, 0, 70}, V(0, 1, 1)})
	if err != nil {
		log.Fatal("failed to make second cut mesh: ", err)
	}

	cmds = slicers.SliceByZ(*m2, *thickness, V(0, 1, 1))

	for i := 0; i < len(cmds); i++ {
		buffer.WriteString(";Layer" + strconv.Itoa(i) + "\n")
		buffer.WriteString(cmds[i].ToGCode())
	}

	alphaG = -30
	buffer.WriteString("G52 " + strconv.Itoa(alphaG) + "\n")

	c = Point{*ox, *oy, *oz}
	cv = c.ToVector()

	alpha = math.Pi * float64(alphaG) / 180.0
	// transposed matrix to rotate around X
	//mx = V(1, 0, 0)
	//my = V(0, math.Cos(alpha), math.Sin(alpha))
	//mz = V(0, -math.Sin(alpha), math.Cos(alpha))
	mx = V(math.Cos(alpha), math.Sin(alpha), 0)
	my = V(-math.Sin(alpha), math.Cos(alpha), 0)
	mz = V(0, 0, 1)

	triangles = make([]Triangle, len(m3.Triangles))
	rotatedMesh = NewMesh(triangles)
	for i, t := range m3.Triangles {
		p1 := c.VectorTo(t.P1).Rotate(mx, my, mz).Add(cv).ToPoint()
		p2 := c.VectorTo(t.P2).Rotate(mx, my, mz).Add(cv).ToPoint()
		p3 := c.VectorTo(t.P3).Rotate(mx, my, mz).Add(cv).ToPoint()
		rotatedMesh.Triangles[i].Fill(p1, p2, p3)
	}

	m4, m3, err := slicers.Cut(&rotatedMesh, Plane{Point{0, 0, 70}, V(0, 0, 1)})
	if err != nil {
		log.Fatal("failed to make third cut mesh: ", err)
	}

	cmds = slicers.SliceByZ(*m3, *thickness, V(1, 1, 0))

	for i := 0; i < len(cmds); i++ {
		buffer.WriteString(";Layer" + strconv.Itoa(i) + "\n")
		buffer.WriteString(cmds[i].ToGCode())
	}

	alphaG = -20
	buffer.WriteString("G42 " + strconv.Itoa(alphaG) + "\n")

	c = Point{*ox, *oy, *oz}
	cv = c.ToVector()

	alpha = math.Pi * float64(alphaG) / 180.0
	// transposed matrix to rotate around X
	//mx = V(math.Cos(alpha), math.Sin(alpha), 0)
	//my = V(-math.Sin(alpha), math.Cos(alpha), 0 )
	//mz = V(0, 0, 1)
	mx = V(1, 0, 0)
	my = V(0, math.Cos(alpha), math.Sin(alpha))
	mz = V(0, -math.Sin(alpha), math.Cos(alpha))

	triangles = make([]Triangle, len(m4.Triangles))
	rotatedMesh = NewMesh(triangles)
	for i, t := range m4.Triangles {
		p1 := c.VectorTo(t.P1).Rotate(mx, my, mz).Add(cv).ToPoint()
		p2 := c.VectorTo(t.P2).Rotate(mx, my, mz).Add(cv).ToPoint()
		p3 := c.VectorTo(t.P3).Rotate(mx, my, mz).Add(cv).ToPoint()
		rotatedMesh.Triangles[i].Fill(p1, p2, p3)
	}

	cmds = slicers.SliceByZ(rotatedMesh, *thickness, V(1, 1, 0))

	for i := 0; i < len(cmds); i++ {
		buffer.WriteString(";Layer" + strconv.Itoa(i) + "\n")
		buffer.WriteString(cmds[i].ToGCode())
	}

	err = ioutil.WriteFile(*gcode, buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save gcode in file: ", err)
	}

}
func main2() {
	kingpin.Parse()

	mesh, err := LoadSTL(*stl)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}

	m2, m1, err := slicers.Cut(mesh, Plane{Point{0, 0, 40}, V(0, 0, 1)})
	if err != nil {
		log.Fatal("failed to cut mesh: ", err)
	}

	cmds := slicers.SliceByZ(*m1, *thickness, V(0, 0, 1))
	//cmds := slicers.SliceWithSlope(*mesh, *thickness, *alpha)

	var buffer bytes.Buffer
	for i := 0; i < len(cmds); i++ {
		buffer.WriteString(";Layer" + strconv.Itoa(i) + "\n")
		buffer.WriteString(cmds[i].ToGCode())
	}

	cmds = slicers.SliceByZ(*m2, *thickness, V(0, 0, 1))
	//cmds := slicers.SliceWithSlope(*mesh, *thickness, *alpha)

	for i := 0; i < len(cmds); i++ {
		buffer.WriteString(";Layer" + strconv.Itoa(i) + "\n")
		buffer.WriteString(cmds[i].ToGCode())
	}

	err = ioutil.WriteFile(*gcode, buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save gcode in file: ", err)
	}
}
