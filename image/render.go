package image

import (
	_ "embed"
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"io"
	"text/template"
)

//go:embed templates/tile.gosvg
var templateString string

var svgTemplate = template.Must(template.New("svgTemplate").Parse(templateString))

type svgTile struct {
	Size        int
	Tracks      []svgTrack
	SwitchPaths []string
}

type svgTrack struct {
	Y      int
	Length int
}

func Render(writer io.Writer, tile domain.Tile, size int) error {
	pixelTracks := make([]svgTrack, 0, len(tile.Tracks))
	offset := len(tile.Tracks) / 2
	y := int(float64(size)/2 - float64(offset*size)*0.1)
	for _, _ = range tile.Tracks {
		generated := svgTrack{Y: y, Length: size}
		pixelTracks = append(pixelTracks, generated)
		y = y + int(float64(size)*0.1)
	}
	switches := computeSwitches(tile.Tracks, pixelTracks, size)
	return svgTemplate.Execute(writer, svgTile{Size: size, Tracks: pixelTracks, SwitchPaths: switches})
}

func computeSwitches(tracks []domain.Track, pxTracks []svgTrack, size int) []string {
	switchWidth := int(0.5 * float64(size))
	switchPaths := make([]string, 0, 0)
	switchStart := size - switchWidth
	switchCp := size - switchWidth/2
	for i, track := range tracks {
		y := pxTracks[i].Y
		for _, sw := range track.Switches {
			target := y + int(float64(sw*size)*0.1)
			merging := i < len(tracks)/2 && sw > 0 ||
				i > len(tracks)/2 && sw < -0
			if merging {
				pxTracks[i].Length = size - switchWidth
			}
			path := fmt.Sprintf("M%d,%d C%d,%d %d,%d, %d,%d", switchStart, y, switchCp, y, switchCp, target, size, target)
			switchPaths = append(switchPaths, path)
		}
	}
	return switchPaths
}
