package renderer

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
	Seed        string
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

type Renderer struct {
	writer io.Writer
	size   int
}

func New(writer io.Writer, size int) Renderer {
	return Renderer{writer: writer, size: size}
}

func (r *Renderer) Render(tile domain.Tile) error {
	pixelTracks := make([]svgTrack, 0, len(tile.Tracks))
	offset := len(tile.Tracks) / 2
	y := int(float64(r.size)/2 - float64(offset*r.size)*trackDistanceRel)
	for _, _ = range tile.Tracks {
		generated := svgTrack{Y: y, Length: r.size}
		pixelTracks = append(pixelTracks, generated)
		y = y + int(float64(r.size)*trackDistanceRel)
	}
	switches := r.computeSwitches(tile.Tracks, pixelTracks)
	bumpers := r.generateBumpers(tile.Tracks, pixelTracks)

	return svgTemplate.Execute(r.writer, svgTile{Size: r.size, Seed: tile.Seed, Tracks: pixelTracks, SwitchPaths: switches, Bumpers: bumpers})
}

func (r *Renderer) computeSwitches(tracks []domain.Track, pxTracks []svgTrack) []string {
	switchWidth := int(switchWidthRel * float64(r.size))
	switchPaths := make([]string, 0, 0)
	switchStart := r.size - switchWidth
	switchCp := r.size - switchWidth/2
	for i, track := range tracks {
		y := pxTracks[i].Y
		for _, sw := range track.Switches {
			target := y + int(float64(sw*r.size)*trackDistanceRel)
			merging := i < len(tracks)/2 && sw > 0 ||
				i > len(tracks)/2 && sw < -0
			if merging {
				pxTracks[i].Length = r.size - switchWidth
			}
			path := fmt.Sprintf("M%d,%d C%d,%d %d,%d, %d,%d", switchStart, y, switchCp, y, switchCp, target, r.size, target)
			switchPaths = append(switchPaths, path)
		}
	}
	return switchPaths
}

func (r *Renderer) generateBumpers(tracks []domain.Track, pxTracks []svgTrack) []svgBumper {
	bumpers := make([]svgBumper, 0, 0)
	for i, track := range tracks {
		bumperLocation := int(switchWidthRel * float64(r.size))
		if track.BumperLeft {
			bumpers = append(bumpers, r.newBumper(bumperLocation, pxTracks[i].Y))
			pxTracks[i].X = bumperLocation
			pxTracks[i].Length = r.size - pxTracks[i].X
		} else if track.BumperRight {
			bumpers = append(bumpers, r.newBumper(r.size-bumperLocation, pxTracks[i].Y))
			pxTracks[i].Length = pxTracks[i].Length - bumperLocation
		}
	}
	return bumpers
}

func (r *Renderer) newBumper(x int, y int) svgBumper {
	bumperSize := int(float64(r.size) * bumperSize)
	return svgBumper{
		X:      x,
		Y:      y - bumperSize/2,
		Width:  bumperSize,
		Height: bumperSize,
	}
}
