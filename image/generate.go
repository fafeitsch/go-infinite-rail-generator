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
	Tracks []int
}

func Render(writer io.Writer, track domain.Tile, size int) error {
	trackSize := size / track.Tracks
	pixelTracks := make([]int, 0, track.Tracks)
	for i := 0; i < track.Tracks; i++ {
		pixelTracks = append(pixelTracks, trackSize*i+trackSize/2)
	}
	return svgTemplate.Execute(writer, svgTile{Size: size, Tracks: pixelTracks})
}
