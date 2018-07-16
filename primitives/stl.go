package primitives

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"
	"os"
	"strings"
	"strconv"
)

type stlHeader struct {
	_     [80]uint8
	Count uint32
}
type stlTriangle struct {
	N, V1, V2, V3 [3]float32
	_             uint16
}

func LoadSTL(path string) (*Mesh, error) {

	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// get file size
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	header := stlHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	expectedSize := int64(header.Count)*50 + 84

	// parse ascii or binary stl
	if info.Size() == expectedSize {
		return loadSTLBinary(file,int(header.Count))
	}
	// rewind to start of file
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	return loadSTLAscii(file)

}

func loadSTLAscii(file *os.File) (*Mesh, error) {
	var triangles []Triangle
	var p1, p2, p3, px Point
	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) < 12 || line[0] != 'v' {
			continue
		}
		fields := strings.Fields(line[7:])
		if len(fields) != 3 {
			continue
		}
		px = Point{parseFloat(fields[0]), parseFloat(fields[1]), parseFloat(fields[2])}
		switch i % 3 {
		case 0:
			p1 = px
		case 1:
			p2 = px
		case 2:
			p3 = px
			triangles = append(triangles, NewTriangle(p1, p2, p3))
		}
		i++
	}
	mesh := NewMesh(triangles)
	return &mesh, scanner.Err()
}
func parseFloat(f string) float64 {
	v, _ := strconv.ParseFloat(f, 32)
	return float64(v)
}

func makeFloat(b []byte) float64 {
	return float64(math.Float32frombits(binary.LittleEndian.Uint32(b)))
}

func loadSTLBinary(file *os.File, count int) (*Mesh, error) {
	b := make([]byte, count*50)
	_, err := io.ReadFull(file, b)
	if err != nil {
		return nil, err
	}

	triangles := make([]Triangle, count)

	work := func(wi, wn int) {
		for i := wi; i < count; i += wn {
			j := i * 50
			p1 := Point{makeFloat(b[j+12: j+16]), makeFloat(b[j+16: j+20]), makeFloat(b[j+20: j+24])}
			p2 := Point{makeFloat(b[j+24: j+28]), makeFloat(b[j+28: j+32]), makeFloat(b[j+32: j+36])}
			p3 := Point{makeFloat(b[j+36: j+40]), makeFloat(b[j+40: j+44]), makeFloat(b[j+44: j+48])}
			triangles[i].Fill(p1, p2, p3)
		}
	}

	DoInParallelAndWait(work)

	mesh := NewMesh(triangles)
	return &mesh, nil
}

func SaveSTL(path string, mesh *Mesh) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	header := stlHeader{}
	header.Count = uint32(len(mesh.Triangles))
	if err := binary.Write(w, binary.LittleEndian, &header); err != nil {
		return err
	}
	for _, t := range mesh.Triangles {
		d := stlTriangle{}
		d.N[0] = float32(t.N.X)
		d.N[1] = float32(t.N.Y)
		d.N[2] = float32(t.N.Z)
		d.V1[0] = float32(t.P1.X)
		d.V1[1] = float32(t.P1.Y)
		d.V1[2] = float32(t.P1.Z)
		d.V2[0] = float32(t.P2.X)
		d.V2[1] = float32(t.P2.Y)
		d.V2[2] = float32(t.P2.Z)
		d.V3[0] = float32(t.P3.X)
		d.V3[1] = float32(t.P3.Y)
		d.V3[2] = float32(t.P3.Z)
		if err := binary.Write(w, binary.LittleEndian, &d); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}
