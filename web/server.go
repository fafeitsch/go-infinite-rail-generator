package web

import (
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"github.com/fafeitsch/go-infinite-rail-generator/image"
	"github.com/fafeitsch/go-infinite-rail-generator/noise"
	"net/http"
	"strconv"
	"strings"
)

func ApiHandler(seed string) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		splitted := strings.Split(r.URL.Path, "/")
		x := splitted[2]
		y := splitted[3]
		if y != "131072.png" {
			return
		}
		fmt.Printf("X is %s\n", x)
		seedString := r.URL.Query().Get("seed")
		if seedString == "" {
			seedString = seed
		}
		hectometer, err := strconv.Atoi(x)
		if err != nil {
			http.Error(writer, fmt.Sprintf("The hectometer query parameter \"%s\" is not a valid number.", r.URL.Query().Get("hectometer")), http.StatusBadRequest)
			return
		}
		numberOfTracks := noise.New(seedString).NumberOfTracks(hectometer)
		writer.Header().Set("Content-Type", "image/svg+xml")
		_ = image.Render(writer, domain.Tile{Tracks: numberOfTracks}, 200)
	}
}
