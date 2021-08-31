package renderer

import (
	_ "embed"
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
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
	"track": func(size int, tracks domain.Tracks) []string {
		beta := size / 3
		gamma := 2 * size / 3
		result := make([]string, 0)
		result = computePaths(float64(size), tracks.Alpha, 0, result)
		result = computePaths(float64(size), tracks.Beta, float64(beta), result)
		result = computePaths(float64(size), tracks.Gamma, float64(gamma), result)
		result = computeBumpers(float64(size), tracks, result)
		return result
	},
}

func computePaths(size float64, column [16]domain.Connectors, x float64, result []string) []string {
	offset := size / 16 / 2
	next := int(math.Ceil(x + size/3/2))
	for i, connectors := range column {
		y := int(math.Floor(float64(i)*size/16 + offset))
		for _, connector := range connectors {
			var path string
			target := int(math.Floor(float64(connector.Slot)*size/16 + offset))
			path = fmt.Sprintf("M %d,%d C%d,%d %d,%d, %d,%d", int(x), y, next, y, next, target, int(x+(math.Ceil(size/3))), target)
			result = append(result, path)
		}
	}
	return result
}

func computeBumpers(size float64, tracks domain.Tracks, result []string) []string {
	gamma := 2 * size / 3
	offset := size / 16 / 2
	bumperWidth := int(size / 16 / 2)
	beta := size/3 - float64(bumperWidth)
	connectorsToGamma := tracks.BuildConnectorMap(domain.Beta, domain.Gamma)
	for i, connectors := range tracks.Gamma {
		if len(connectors) == 0 && connectorsToGamma[i] {
			y := int(math.Floor(float64(i)*size/16 + offset))
			path := fmt.Sprintf("M %d %d h %d v %d h %d Z", int(gamma), y-bumperWidth/2, bumperWidth, bumperWidth, -bumperWidth)
			result = append(result, path)
		}
	}
	connectorsToBeta := tracks.BuildConnectorMap(domain.Alpha, domain.Beta)
	for i, connectors := range tracks.Beta {
		if !connectorsToBeta[i] && len(connectors) > 0 {
			y := int(math.Floor(float64(i)*size/16 + offset))
			path := fmt.Sprintf("M %d %d h %d v %d h %d Z", int(beta), y-bumperWidth/2, bumperWidth, bumperWidth, -bumperWidth)
			result = append(result, path)
		}
	}
	return result
}

var svgTemplate = template.Must(template.New("svgTemplate").Funcs(templateFunctions).Parse(templateString))

type svgTile struct {
	Size  int
	Rails domain.Tracks
}

type Renderer struct {
	writer io.Writer
	size   int
}

func New(writer io.Writer, size int) Renderer {
	return Renderer{writer: writer, size: size}
}

func (r *Renderer) Render(tile *domain.Tile) error {
	return svgTemplate.Execute(r.writer, svgTile{Size: r.size, Rails: tile.Tracks})
}
