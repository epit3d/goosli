package slicers

import (
	"github.com/l1va/goosli"
	"bytes"
	"strconv"
	"io/ioutil"
	"log"
	"fmt"
)

// Slice - Slicing on layers by simple algo
func Slice5aByCenter(mesh *goosli.Mesh, thickness float64) bytes.Buffer {

	epsilon := 10.0 //TODO: move to params

	layers := SliceByZ(mesh, thickness, goosli.V(0, 0, 1))
	layersToGcode(layers, "/home/l1va/debug.gcode")

	centers := calculateCenters(layers)
	centersToDebugFile(centers, "/home/l1va/debug.txt")

	simplified := SimplifyLine(centers, epsilon)
	centersToDebugFile(simplified, "/home/l1va/debug_simplified.txt")

	var buffer bytes.Buffer
	up := mesh
	down := mesh
	err := fmt.Errorf("temp") //TODO: fixme
	layersCount := 0
	absangleZ := 0
	absangleX := 0


	for i := 1; i < len(simplified); i++ {
		v := simplified[i-1].VectorTo(simplified[i])
		if i < len(simplified)-1 {
			up, down, err = Cut(up, goosli.Plane{simplified[i], v}) //TODO: rotate points too
			if err != nil {
				log.Fatal("failed to cut mesh by plane: ", err)
			}
		} else {
			down = up
		}
		//TODO: fixme
		absangleZ = calcZ(simplified[i-1], simplified[i], nil)
		absangleX = calcX(simplified[i-1], simplified[i], nil)
		println(absangleX," ", absangleZ)

		down = rotateXZ(absangleX,absangleZ, down, &buffer, goosli.Point{0,0,0})

		layersCount += slicePart(down, v, thickness, layersCount, &buffer)
	}

	return buffer
}

func layersToGcode(layers []goosli.Layer, filename string) {
	var bf bytes.Buffer
	for i := 0; i < len(layers); i++ {
		bf.WriteString(";Layer" + strconv.Itoa(i) + "\n")
		bf.WriteString(layers[i].ToGCode())
	}
	err := ioutil.WriteFile(filename, bf.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save upper mesh: ", err)
	}
}

func calculateCenters(layers []goosli.Layer) []goosli.Point {
	var centers []goosli.Point
	for _, layer := range layers {
		x, y, z, count := 0.0, 0.0, 0.0, 0
		for _, path := range layer.Paths {
			for _, line := range path.Lines {
				x += line.P1.X + line.P2.X
				y += line.P1.Y + line.P2.Y
				z += line.P1.Z + line.P2.Z
			}
			count += len(path.Lines) * 2
		}
		if count > 0 {
			countF := float64(count)
			centers = append(centers, goosli.Point{x / countF, y / countF, z / countF})
		}
	}
	return centers
}

func centersToDebugFile(ps []goosli.Point, filename string) {
	var b bytes.Buffer
	for i := 0; i < len(ps)-1; i++ {
		b.WriteString("line ")
		b.WriteString(ps[i].ToString2())
		b.WriteString(ps[i+1].ToString2() + "\n")
	}

	err := ioutil.WriteFile(filename, b.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save debug in file: ", err)
	}
}
