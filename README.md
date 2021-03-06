# goosli

[![Build Status](https://travis-ci.org/l1va/goosli.svg?branch=master)](https://travis-ci.org/l1va/goosli)

Package to implement your own slicer using existing possibilities: cutting, crossing, 
slicing by vector, mesh simplification, line simplification etc.

For example you can slice by profile (for <b>5axes 3d printer</b>) - find center of 
each layer, make line from points, simplify line - and slice by this line. 
See slicers/slice_by_profile.go. For sure your 3d-printer have to support bed rotations.

<i>Feel free to open issues or implement your slicing algorithms.</i>

# goosli-colorizer
Save true/false to file for each triangle in stl according to logic (now it is big bias from Z axe)

# goosli-cutter
Cut stl in two stls by required plane.

# goosli-simplifier
Simplifies stl to required count of triangles.

### Viewer
[Spycer - https://github.com/l1va/spycer](https://github.com/l1va/spycer)

### Get binaries
Do not forget to place <b>data directory</b> near your binary. 
##### Linux (tested on Ubuntu 16.04)
```bash
cd github.com/l1va/goosli
go install cmd/goosli/goosli.go 
go install cmd/goosli_colorizer/goosli_colorizer.go
go install cmd/goosli_cutter/goosli_cutter.go 
go install cmd/goosli_simplifier/goosli_simplifier.go
```
##### Windows (tested on Windows 10)
```bash
cd github.com/l1va/goosli
GOOS=windows GOARCH=amd64 go build -o goosli cmd/goosli/goosli.go
GOOS=windows GOARCH=amd64 go build -o goosli_analyzer cmd/goosli_analyzer/goosli_analyzer.go 
GOOS=windows GOARCH=amd64 go build -o goosli_colorizer cmd/goosli_colorizer/goosli_colorizer.go
GOOS=windows GOARCH=amd64 go build -o goosli_cutter cmd/goosli_cutter/goosli_cutter.go 
GOOS=windows GOARCH=amd64 go build -o goosli_simplifier cmd/goosli_simplifier/goosli_simplifier.go 
```

### Technical moments
Rotations are always about global axes, if your rotation axis does not match with global 
axis - you can use PlaneCenter parameter and shift your plane to match.

### Thanks
A lot of ideas and code was taken from various [fogleman](https://github.com/fogleman) 
repos. Thank you!
