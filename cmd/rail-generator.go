package main

import (
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"github.com/fafeitsch/go-infinite-rail-generator/image"
	"github.com/fafeitsch/go-infinite-rail-generator/noise"
	"log"
	"os"
)

func main() {
	tracks := noise.New(42).NumberOfTracks(66)
	tile := domain.Tile{Tracks: tracks}
	err := image.Render(tile, 500, os.Stdout)
	if err != nil {
		log.Fatalf("could not render image: %v", err)
	}
}
