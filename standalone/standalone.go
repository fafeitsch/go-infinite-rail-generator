package standalone

import (
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/generator"
	"github.com/fafeitsch/go-infinite-rail-generator/renderer"
	"os"
)

type RenderOptions struct {
	Seed       string
	TownNames  []string
	Hectometer int
	Size       int
}

func RenderSingleTile(options RenderOptions) error {
	gen := generator.New(options.Seed)
	gen.TownNames = options.TownNames
	size := options.Size
	if size == 0 {
		size = 200
	}
	tile := gen.Generate(options.Hectometer)
	rn := renderer.New(os.Stdout, size)
	err := rn.Render(tile)
	if err != nil {
		return fmt.Errorf("could not render tile %d: %v", options.Hectometer, err)
	}
	return nil
}
