package renderer

import (
	_ "embed"
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/world"
	"io"
	"math"
	"text/template"
)

//go:embed templates/tile.gosvg
var templateString string

var templateFunctions = template.FuncMap{
	"add": func(a int, b int) int {
		return a + b
	},
	"half": func(a int) int {
		return a / 2
	},
	"track": func(size int, tracks world.Tracks) []string {
		result := computePaths(float64(size), tracks)
		// not needed for now – maybe bumbers will be set explicitly
		// result = computeBumpers(float64(size), tracks, result)
		return result
	},
	"trackLabel": computeTrackLabels,
	"platform":   computePlatforms,
}

func computePaths(size float64, tracks world.Tracks) []string {
	result := make([]string, 0, len(tracks))
	offset := size / 16 / 2
	for _, connector := range tracks {
		x := float64(connector.SourceColumn) * size / 3
		y := float64(connector.SourceTrack)*size/16 + offset
		var path string
		xTarget := float64(connector.TargetColumn) * size / 3
		controlPointX := (xTarget + x) / 2
		yTarget := float64(connector.TargetTrack)*size/16 + offset
		path = fmt.Sprintf("d=\"M %f,%f C%f,%f %f,%f, %f,%f\"", x, y, controlPointX, y, controlPointX, yTarget, xTarget, yTarget)
		if connector.Junction {
			path = fmt.Sprintf("%s stroke-dasharray=\"2\"", path)
		}
		result = append(result, path)
	}
	return result
}

func computeBumpers(size float64, tracks world.Tracks, result []string) []string {
	// gamma := 2 * size / 3
	// offset := size / 16 / 2
	// bumperWidth := int(size / 16 / 2)
	// beta := size/3 - float64(bumperWidth)
	// connectorsToGamma := tracks.BuildConnectorMap(world.Beta, world.Gamma)
	// for i, connectors := range tracks.Get(world.Gamma) {
	// 	if len(connectors) == 0 && connectorsToGamma[i] {
	// 		y := int(math.Floor(float64(i)*size/16 + offset))
	// 		path := fmt.Sprintf("d=\"M %d %d h %d v %d h %d Z\"", int(gamma), y-bumperWidth/2, bumperWidth, bumperWidth, -bumperWidth)
	// 		result = append(result, path)
	// 	}
	// }
	// connectorsToBeta := tracks.BuildConnectorMap(world.Alpha, world.Beta)
	// for i, connectors := range tracks.Get(world.Beta) {
	// 	if !connectorsToBeta[i] && len(connectors) > 0 {
	// 		y := int(math.Floor(float64(i)*size/16 + offset))
	// 		path := fmt.Sprintf("d=\"M %d %d h %d v %d h %d Z\"", int(beta), y-bumperWidth/2, bumperWidth, bumperWidth, -bumperWidth)
	// 		result = append(result, path)
	// 	}
	// }
	// return result
	return []string{}
}

func computePlatforms(size int, platforms []world.Platform) []string {
	result := make([]string, 0, 0)
	offset := float64(size) / 16 / 2
	for _, platform := range platforms {
		y := float64(platform.Track)*float64(size)/16 - offset + 2
		x := float64(platform.Column) * float64(size) / 3
		height := size/16 - 4
		result = append(result, fmt.Sprintf("M %f %f h %f v %d h -%f Z", x, y, float64(size)/3, height, float64(size)/3))
	}
	return result
}

func computeTrackLabels(size int, tracks world.Tracks) []string {
	result := make([]string, 0, 0)
	for _, connector := range tracks {
		y := math.Floor(float64(connector.SourceTrack) * float64(size) / 16)
		if connector.Label != "" {
			x := float64(size) * float64(connector.SourceColumn) / 3
			element := fmt.Sprintf("x=\"%f\" y=\"%f\" font-size=\"0.5em\">%s</text>", x, y, connector.Label)
			result = append(result, element)
		}
	}
	return result
}

var svgTemplate = template.Must(template.New("svgTemplate").Funcs(templateFunctions).Parse(templateString))

type svgTile struct {
	*world.Tile
	Size int
}

type Renderer struct {
	writer io.Writer
	size   int
}

func New(writer io.Writer, size int) Renderer {
	return Renderer{writer: writer, size: size}
}

func (r *Renderer) Render(tile world.Tile) error {
	return svgTemplate.Execute(r.writer, svgTile{Size: r.size, Tile: &tile})
}
