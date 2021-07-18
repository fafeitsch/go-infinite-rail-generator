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

var templateFunctions = template.FuncMap{
	"add": func(a int, b int) int {
		return a + b
	},
}

var svgTemplate = template.Must(template.New("svgTemplate").Funcs(templateFunctions).Parse(templateString))

const switchWidthRel = 0.5
const trackDistanceRel = 0.1
const bumperSize = 0.025

type svgTile struct {
	Size        int
	Tracks      []svgTrack
	SwitchPaths []string
	Bumpers     []svgBumper
}

type svgTrack struct {
	X      int
	Y      int
	Length int
}

type svgBumper struct {
	X      int
	Y      int
	Width  int
	Height int
}

func Render(writer io.Writer, tile domain.Tile, size int) error {
	pixelTracks := make([]svgTrack, 0, len(tile.Tracks))
	offset := len(tile.Tracks) / 2
	y := int(float64(size)/2 - float64(offset*size)*trackDistanceRel)
	bumpers := make([]svgBumper, 0, 0)
	for _, track := range tile.Tracks {
		generated := svgTrack{Y: y, Length: size}
		if track.BumperLeft {
			bumpers = append(bumpers, newBumper(int(switchWidthRel*float64(size)), y, size))
			generated.X = int(switchWidthRel * float64(size))
			generated.Length = size - generated.X
		} else if track.BumperRight {
			bumpers = append(bumpers, newBumper(size-int(switchWidthRel*float64(size)), y, size))
			generated.Length = generated.Length - int(switchWidthRel*float64(size))
		}
		pixelTracks = append(pixelTracks, generated)
		y = y + int(float64(size)*trackDistanceRel)
	}
	switches := computeSwitches(tile.Tracks, pixelTracks, size)

	return svgTemplate.Execute(writer, svgTile{Size: size, Tracks: pixelTracks, SwitchPaths: switches, Bumpers: bumpers})
}

func computeSwitches(tracks []domain.Track, pxTracks []svgTrack, size int) []string {
	switchWidth := int(switchWidthRel * float64(size))
	switchPaths := make([]string, 0, 0)
	switchStart := size - switchWidth
	switchCp := size - switchWidth/2
	for i, track := range tracks {
		y := pxTracks[i].Y
		for _, sw := range track.Switches {
			target := y + int(float64(sw*size)*trackDistanceRel)
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

func newBumper(x int, y int, size int) svgBumper {
	bumperSize := int(float64(size) * bumperSize)
	return svgBumper{
		X:      x,
		Y:      y - bumperSize/2,
		Width:  bumperSize,
		Height: bumperSize,
	}
}
