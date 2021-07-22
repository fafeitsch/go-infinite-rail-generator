package web

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/noise"
	"github.com/fafeitsch/go-infinite-rail-generator/renderer"
	"github.com/fafeitsch/go-infinite-rail-generator/version"
	"net/http"
	"path"
	"strconv"
	"strings"
)

func ApiHandler(defaultNoise *noise.Noise, shift int) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		var head string
		originalPath := r.URL.Path
		head, r.URL.Path = shiftPath(r.URL.Path)
		switch head {
		case "tiles":
			serveTile(defaultNoise, shift, writer, r)
			break
		case "config":
			serveConfig(defaultNoise, writer, r)
			break
		default:
			r.URL.Path = originalPath
			http.FileServer(http.Dir("./html")).ServeHTTP(writer, r)
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

func serveTile(defaultNoise *noise.Noise, shift int, writer http.ResponseWriter, r *http.Request) {
	y := r.URL.Query().Get("vertical")
	if y != "131072" {
		return
	}
	seedString := r.URL.Query().Get("seed")
	hectometer, err := strconv.Atoi(r.URL.Query().Get("tile"))
	if err != nil {
		http.Error(writer, fmt.Sprintf("The tile query parameter \"%s\" is not a valid number.", r.URL.Query().Get("tile")), http.StatusBadRequest)
		return
	}
	var aNoise *noise.Noise
	if seedString == "" || seedString == defaultNoise.Seed {
		aNoise = defaultNoise
	} else {
		aNoise = noise.New(seedString)
	}
	tile := aNoise.Generate(hectometer + shift)
	writer.Header().Set("Content-Type", "image/svg+xml")
	rn := renderer.New(writer, 200)
	_ = rn.Render(tile)
}

type configDto struct {
	Version     string `json:"version"`
	BuildTime   string `json:"buildTime"`
	DefaultSeed string `json:"defaultSeed"`
}

func serveConfig(defaultNoise *noise.Noise, writer http.ResponseWriter, r *http.Request) {
	config := configDto{
		DefaultSeed: defaultNoise.Seed,
		BuildTime:   version.BuildTime,
		Version:     version.BuildVersion,
	}
	_ = json.NewEncoder(writer).Encode(config)
}
