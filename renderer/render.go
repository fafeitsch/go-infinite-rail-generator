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
		beta := size / 3
		gamma := 2 * size / 3
		result := make([]string, 0)
		result = computePaths(float64(size), tracks.Alpha, 0, result)
		result = computePaths(float64(size), tracks.Beta, float64(beta), result)
		result = computePaths(float64(size), tracks.Gamma, float64(gamma), result)
		// not needed for now – maybe bumbers will be set explicitly
		// result = computeBumpers(float64(size), tracks, result)
		return result
	},
	"trackLabel": computeTrackLabels,
	"platform":   computePlatforms,
}

func computePaths(size float64, column [16]world.Connectors, x float64, result []string) []string {
	offset := size / 16 / 2
	for i, connectors := range column {
		y := int(math.Floor(float64(i)*size/16 + offset))
		for _, connector := range connectors {
			var path string
			xDash := size / 3
			if connector.Target == world.Gamma {
				xDash = 2 * size / 3
			} else if connector.Target == world.Omega {
				xDash = size
			}
			next := int(x) + int(math.Ceil((xDash-x)/2))
			yTarget := int(math.Floor(float64(connector.Track)*size/16 + offset))
			path = fmt.Sprintf("d=\"M %d,%d C%d,%d %d,%d, %d,%d\"", int(x), y, next, y, next, yTarget, int(x+(math.Ceil(xDash-x))), yTarget)
			if connector.Junction {
				path = fmt.Sprintf("%s stroke-dasharray=\"2\"", path)
			}
			result = append(result, path)
		}
	}
	return result
}

func computeBumpers(size float64, tracks world.Tracks, result []string) []string {
	gamma := 2 * size / 3
	offset := size / 16 / 2
	bumperWidth := int(size / 16 / 2)
	beta := size/3 - float64(bumperWidth)
	connectorsToGamma := tracks.BuildConnectorMap(world.Beta, world.Gamma)
	for i, connectors := range tracks.Gamma {
		if len(connectors) == 0 && connectorsToGamma[i] {
			y := int(math.Floor(float64(i)*size/16 + offset))
			path := fmt.Sprintf("d=\"M %d %d h %d v %d h %d Z\"", int(gamma), y-bumperWidth/2, bumperWidth, bumperWidth, -bumperWidth)
			result = append(result, path)
		}
	}
	connectorsToBeta := tracks.BuildConnectorMap(world.Alpha, world.Beta)
	for i, connectors := range tracks.Beta {
		if !connectorsToBeta[i] && len(connectors) > 0 {
			y := int(math.Floor(float64(i)*size/16 + offset))
			path := fmt.Sprintf("d=\"M %d %d h %d v %d h %d Z\"", int(beta), y-bumperWidth/2, bumperWidth, bumperWidth, -bumperWidth)
			result = append(result, path)
		}
	}
	return result
}

func computePlatforms(size int, platforms world.Platforms) []string {
	result := make([]string, 0, 0)
	offset := size / 16 / 2
	for i, _ := range platforms.Alpha {
		y := int(math.Floor(float64(i)*float64(size)/16)) - offset + 2
		height := size/16 - 4
		if platforms.Alpha[i] {
			result = append(result, fmt.Sprintf("M 0 %d h %d v %d H 0 Z", y, int(math.Ceil(float64(size)/3)), height))
		}
		if platforms.Beta[i] {
			result = append(result, fmt.Sprintf("M %d %d h %d v %d h -%d Z", size/3, y, size/3, height, size/3))
		}
		if platforms.Gamma[i] {
			result = append(result, fmt.Sprintf("M %d %d H %d v %d h -%d Z", 2*size/3, y, size, height, int(math.Ceil(2*float64(size)/3))))
		}
	}
	return result
}

func computeTrackLabels(size int, tracks world.Tracks) []string {
	result := make([]string, 0, 0)
	for i, _ := range tracks.Alpha {
		y := math.Floor(float64(i) * float64(size) / 16)
		for _, connector := range tracks.Alpha[i] {
			if connector.Label != "" {
				element := fmt.Sprintf("x=\"%f\" y=\"%f\" font-size=\"0.5em\">%s</text>", float64(size)*1/3/2, y, connector.Label)
				result = append(result, element)
			}
		}
		for _, connector := range tracks.Beta[i] {
			if connector.Label != "" {
				element := fmt.Sprintf("x=\"%f\" y=\"%f\" font-size=\"0.5em\">%s</text>", float64(size)*1.5/3, y, connector.Label)
				result = append(result, element)
			}
		}
		for _, connector := range tracks.Gamma[i] {
			if connector.Label != "" {
				element := fmt.Sprintf("x=\"%f\" y=\"%f\" font-size=\"0.5em\">%s</text>", float64(size)*2.5/3, y, connector.Label)
				result = append(result, element)
			}
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
