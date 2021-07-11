package web

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"github.com/fafeitsch/go-infinite-rail-generator/image"
	"github.com/fafeitsch/go-infinite-rail-generator/noise"
	"net/http"
	"path"
	"strconv"
	"strings"
)

func ApiHandler(seed string) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		var head string
		head, r.URL.Path = shiftPath(r.URL.Path)
		switch head {
		case "tiles":
			serveTile(seed, writer, r)
			break
		case "config":
			serveConfig(seed, writer, r)
			break
		}
	}
}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func serveTile(defaultSeed string, writer http.ResponseWriter, r *http.Request) {
	y := r.URL.Query().Get("vertical")
	if y != "131072" {
		return
	}
	seedString := r.URL.Query().Get("seed")
	if seedString == "" {
		seedString = defaultSeed
	}
	hectometer, err := strconv.Atoi(r.URL.Query().Get("hectometer"))
	if err != nil {
		http.Error(writer, fmt.Sprintf("The hectometer query parameter \"%s\" is not a valid number.", r.URL.Query().Get("hectometer")), http.StatusBadRequest)
		return
	}
	numberOfTracks := noise.New(seedString).NumberOfTracks(hectometer)
	writer.Header().Set("Content-Type", "image/svg+xml")
	_ = image.Render(writer, domain.Tile{Tracks: numberOfTracks}, 200)
}

type configDto struct {
	DefaultSeed string `json:"defaultSeed"`
}

func serveConfig(defaultSeed string, writer http.ResponseWriter, r *http.Request) {
	config := configDto{DefaultSeed: defaultSeed}
	_ = json.NewEncoder(writer).Encode(config)
}
