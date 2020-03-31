#!/bin/bash



cd ~/go/src/github.com/l1va/goosli
go install cmd/goosli/goosli.go
go install cmd/goosli_colorizer/goosli_colorizer.go
go install cmd/goosli_cutter/goosli_cutter.go
go install cmd/goosli_simplifier/goosli_simplifier.go
go install cmd/goosli_analyzer/goosli_analyzer.go

cp ~/go/bin/goosli* ~/projects/spycer/
cp -r ./data ~/projects/spycer

