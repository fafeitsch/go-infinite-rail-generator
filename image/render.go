package image

import (
	_ "embed"
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"io"
	"text/template"
)

//go:embed templates/tile.gosvg
var templateString string

var svgTemplate = template.Must(template.New("svgTemplate").Parse(templateString))

type svgTile struct {
	Size   int
	Tracks []svgTrack
}

type svgTrack struct {
	Y      int
	Length int
}

func Render(writer io.Writer, track domain.Tile, size int) error {
	pixelTracks := make([]svgTrack, 0, len(track.Tracks))
	offset := len(track.Tracks) / 2
	y := int(float64(size)/2 - float64(offset*size)*0.1)
	for i := 0; i < len(track.Tracks); i++ {
		pixelTracks = append(pixelTracks, svgTrack{Y: y, Length: size})
		y = y + int(float64(size)*0.1)
	}
	return svgTemplate.Execute(writer, svgTile{Size: size, Tracks: pixelTracks})
}
