package vip

import (
	"bufio"
	"os"

	. "github.com/l1va/goosli/primitives"
	. "github.com/l1va/goosli/slicers"
)

//PrepareLayers add Walls and fill layers
func PrepareLayers(layers []Layer, settings Settings, planes []Plane, fullPanes []Plane) []Layer {
	addWallsCount := int(settings.WallThickness / settings.GcodeSettings.LineWidth)
	if addWallsCount > 0 {
		for i, layer := range layers { //TODO: in parallel
			println("Layer %v, num path %v", i, len(layer.Paths))
			for _, pt := range layer.Paths {
				if len(pt.Points) < 2 { //TODO: remove this
					continue
				}
				offs := offset(pt, addWallsCount, settings.GcodeSettings.LineWidth, layer.Norm)
				layers[i].MiddlePs = append(layers[i].MiddlePs, offs[:len(offs)-1]...)
				layers[i].InnerPs = append(layers[i].InnerPs, offs[len(offs)-1])
			}
		}
	}
	return FillLayers(layers, planes, fullPanes, settings)
}

func SkirtPathes(first Layer, count int, lineWidth float64) Layer {
	if count == 0 {
		return first
	}
	var res []Path
	for _, pt := range first.Paths {
		if len(pt.Points) < 2 { //TODO: remove this
			continue
		}
		offs := offset(pt, count, lineWidth, first.Norm.Reverse())
		for j := range offs {
			res = append(res, offs[len(offs)-1-j])
		}
		res = append(res, pt)
	}
	first.Paths = res
	return first
}
func offset(pth Path, addWallsCount int, nozzle float64, norm Vector) []Path {
	var res []Path
	res = append(res, MakeOffset(pth, nozzle, norm))
	for i := 1; i < addWallsCount; i++ {
		res = append(res, MakeOffset(res[i-1], nozzle, norm))
	}
	return res

}

func readPlanes(f string) []AnalyzedPlane {
	file, err := os.Open(f)
	if err != nil {
		println("File reading error: " + err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ans := []AnalyzedPlane{}
	for scanner.Scan() {
		ans = append(ans, ParsePlane(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		println("Errors during reading: " + err.Error())
	}
	return ans
}
