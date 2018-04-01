package goosli

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"
	"os"
	"runtime"
	"sync"
)

type stlHeader struct {
	_     [80]uint8
	Count uint32
}

func LoadSTL(path string) (*Mesh, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return loadSTLBinary(file)

}

func makeFloat(b []byte) float64 {
	return float64(math.Float32frombits(binary.LittleEndian.Uint32(b)))
}

func loadSTLBinary(file *os.File) (*Mesh, error) {
	reader := bufio.NewReader(file)
	header := stlHeader{}
	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	count := int(header.Count)
	b := make([]byte, count*50)
	_, err := io.ReadFull(reader, b)
	if err != nil {
		return nil, err
	}

	triangles := make([]Triangle, count)

	work := func(wi, wn int) {
		for i := wi; i < count; i += wn {
			j := i * 50
			v1 := Vector{makeFloat(b[j+12: j+16]), makeFloat(b[j+16: j+20]), makeFloat(b[j+20: j+24])}
			v2 := Vector{makeFloat(b[j+24: j+28]), makeFloat(b[j+28: j+32]), makeFloat(b[j+32: j+36])}
			v3 := Vector{makeFloat(b[j+36: j+40]), makeFloat(b[j+40: j+44]), makeFloat(b[j+44: j+48])}
			triangles[i].fill(v1, v2, v3)
		}
	}

	doParallel(work)

	mesh := NewMesh(triangles)
	return &mesh, nil
}

func doParallel(work func(wi, wn int)) {
	wn := runtime.NumCPU()
	var wg sync.WaitGroup
	for wi := 0; wi < wn; wi++ {
		wg.Add(1)
		go func(wi, wn int) {
			work(wi, wn)
			wg.Done()
		}(wi, wn)
	}
	wg.Wait()
}
