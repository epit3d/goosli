package main

import (
	"log"
	"github.com/l1va/goosli/slicers"
	"bytes"
	"io/ioutil"
	. "github.com/l1va/goosli"
)

func main() {
	path := "examples/bowser.stl"
	gcodePath := "examples/bowser.gcode"

	mesh, err := LoadSTL(path)
	if err != nil {
		log.Fatal("failed to load mesh: ", err)
	}
	cmds := slicers.Slice3DOF(*mesh)

	var buffer bytes.Buffer
	for i := 0; i < len(cmds); i++ {
		buffer.WriteString(cmds[i].ToGCode() + "\n")
	}

	err = ioutil.WriteFile(gcodePath, buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal("failed to save gcode in file: ", err)
	}
}
