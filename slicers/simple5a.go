package slicers

import (
	"github.com/l1va/goosli"
	"math"
	"bytes"
	"strconv"
	"fmt"
	"log"
	"io/ioutil"
)
// WARNING: not working
// Slice - Slicing on layers by simple algo
func Slice5a(mesh *goosli.Mesh, thickness float64) bytes.Buffer {

	sinValue := 0.90
	c := goosli.Point{0, 0, 0} //TODO:

	var b bytes.Buffer
	layers := 0
	z := goosli.V(0, 0, 1)
	found := true
	var store = map[int]int{}
	for found {
		found = false

		var mint *goosli.Triangle
		minz := math.MaxFloat64
		debugArr := []goosli.Triangle{}
		for _, t := range mesh.Triangles {
			store[int(findSin(z, &t)*100)] += 1
			if findSin(z, &t) > sinValue {
				debugArr = append(debugArr, t)
				nz := t.MinZ(z)
				if nz < minz {
					minz = nz
					mint = &t
				}
			}
		}

		toFile(debugArr)
		for i := -100; i < 100; i += 1 {
			fmt.Println(i, store[i])
		}
		if mint != nil {
			found = true

			point := z.MulScalar(minz).ToPoint()

			up, down, err := Cut(mesh, goosli.Plane{point, z})
			if err != nil {
				log.Fatal("failed to cut mesh by plane: ", err)
			}

			layers += slicePart(down, z, thickness, layers, &b)

			angleZ := calcZ(c, point, mint)
			angleX := calcX(c, point, mint)

			if angleZ != 0 {
				up = rotateZ(angleZ, up, &b, c)
			}
			mesh = rotateX(angleX, up, &b, c)

		} else {
			slicePart(mesh, z, thickness, layers, &b)
		}
	}
	return b
}

func toFile(ts []goosli.Triangle) {
	var b bytes.Buffer
	for _,t := range(ts) {
		b.WriteString("triangle " )
		b.WriteString(t.P1.ToString2() )
		b.WriteString(t.P2.ToString2() )
		b.WriteString(t.P3.ToString2() +"\n")
	}

	err := ioutil.WriteFile("/home/l1va/debug.txt", b.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save debug in file: ", err)
	}
}

func calcZ(c goosli.Point, p goosli.Point, t *goosli.Triangle) int {
	//TODO: rethink, it is not honest
	v := c.VectorTo(p).ProjectOnPlane(goosli.Plane{goosli.Point{0, 0, 0}, goosli.V(0, 0, 1)})
	return int(v.Angle(goosli.V(1, 0, 0)))
}

func calcX(c goosli.Point, p goosli.Point, t *goosli.Triangle) int {
	//TODO: rethink, it is not honest
	v := c.VectorTo(p).ProjectOnPlane(goosli.Plane{goosli.Point{0, 0, 0}, goosli.V(1, 0, 0)})
	return int(v.Angle(goosli.V(0, 0, 1)))
}

func rotateXZ(angleX int, angleZ int,mesh *goosli.Mesh, b *bytes.Buffer, c goosli.Point) *goosli.Mesh {
	b.WriteString("G62 " + strconv.Itoa(angleX) +" "+strconv.Itoa(angleZ)+ "\n")
	mesh = rotateX(angleX,mesh, b, c)
	return rotateZ(angleZ, mesh, b, c)
}

func rotateX(angle int,mesh *goosli.Mesh, b *bytes.Buffer, c goosli.Point) *goosli.Mesh {
	//b.WriteString("G42 " + strconv.Itoa(absangle) + "\n")

	cv := c.ToVector()

	alpha := math.Pi * float64(angle) / 180.0
	// transposed matrix to rotate around X
	mx := goosli.V(1, 0, 0)
	my := goosli.V(0, math.Cos(alpha), math.Sin(alpha))
	mz := goosli.V(0, -math.Sin(alpha), math.Cos(alpha))

	triangles := make([]goosli.Triangle, len(mesh.Triangles))
	rotatedMesh := goosli.NewMesh(triangles)
	for i, t := range mesh.Triangles {
		p1 := c.VectorTo(t.P1).Rotate(mx, my, mz).Add(cv).ToPoint()
		p2 := c.VectorTo(t.P2).Rotate(mx, my, mz).Add(cv).ToPoint()
		p3 := c.VectorTo(t.P3).Rotate(mx, my, mz).Add(cv).ToPoint()
		rotatedMesh.Triangles[i].Fill(p1, p2, p3)
	}
	return &rotatedMesh
}
func rotateZ(angle int, mesh *goosli.Mesh, b *bytes.Buffer, c goosli.Point) *goosli.Mesh {
	//b.WriteString("G52 " + strconv.Itoa(absangle) + "\n")

	cv := c.ToVector()

	alpha := math.Pi * float64(angle) / 180.0
	// transposed matrix to rotate around Z
	mx := goosli.V(math.Cos(alpha), math.Sin(alpha), 0)
	my := goosli.V(-math.Sin(alpha), math.Cos(alpha), 0)
	mz := goosli.V(0, 0, 1)

	triangles := make([]goosli.Triangle, len(mesh.Triangles))
	rotatedMesh := goosli.NewMesh(triangles)
	for i, t := range mesh.Triangles {
		p1 := c.VectorTo(t.P1).Rotate(mx, my, mz).Add(cv).ToPoint()
		p2 := c.VectorTo(t.P2).Rotate(mx, my, mz).Add(cv).ToPoint()
		p3 := c.VectorTo(t.P3).Rotate(mx, my, mz).Add(cv).ToPoint()
		rotatedMesh.Triangles[i].Fill(p1, p2, p3)
	}
	return &rotatedMesh
}

func slicePart(mesh *goosli.Mesh, v goosli.Vector, thickness float64, start int, b *bytes.Buffer) int {
	cmds := SliceByZ(mesh, thickness, v)

	for i := 0; i < len(cmds); i++ {
		b.WriteString(";Layer" + strconv.Itoa(i+start) + "\n")
		b.WriteString(cmds[i].ToGCode())
	}
	return len(cmds)
}

func findSin(z goosli.Vector, t *goosli.Triangle) float64 {
	return math.Abs(z.Dot(t.N)) / z.Length() / t.N.Length()
}
