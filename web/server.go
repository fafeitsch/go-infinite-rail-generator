package web

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/image"
	"github.com/fafeitsch/go-infinite-rail-generator/noise"
	"net/http"
	"path"
	"strconv"
	"strings"
)

func ApiHandler(defaultNoise *noise.Noise) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		var head string
		originalPath := r.URL.Path
		head, r.URL.Path = shiftPath(r.URL.Path)
		switch head {
		case "tiles":
			serveTile(defaultNoise, writer, r)
			break
		case "config":
			serveConfig(defaultNoise, writer, r)
			break
		default:
			r.URL.Path = originalPath
			http.FileServer(http.Dir("./web/app/dist/app")).ServeHTTP(writer, r)
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

func serveTile(defaultNoise *noise.Noise, writer http.ResponseWriter, r *http.Request) {
	y := r.URL.Query().Get("vertical")
	if y != "131072" {
		return
	}
	seedString := r.URL.Query().Get("seed")
	offset := 0
	if r.URL.Query().Get("offset") != "" {
		var err error
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(writer, fmt.Sprintf("The offset %s is not a valid number.", r.URL.Query().Get("offset")), http.StatusBadRequest)
		}
	}
	hectometer, err := strconv.Atoi(r.URL.Query().Get("hectometer"))
	if err != nil {
		http.Error(writer, fmt.Sprintf("The hectometer query parameter \"%s\" is not a valid number.", r.URL.Query().Get("hectometer")), http.StatusBadRequest)
		return
	}
	var aNoise *noise.Noise
	if seedString == "" || seedString == defaultNoise.Seed {
		aNoise = defaultNoise
	} else {
		aNoise = noise.New(seedString)
	}
	tile := aNoise.Generate(hectometer - offset)
	writer.Header().Set("Content-Type", "image/svg+xml")
	renderer := image.New(writer, 200)
	_ = renderer.Render(tile)
}

type configDto struct {
	DefaultSeed string `json:"defaultSeed"`
}

func serveConfig(defaultNoise *noise.Noise, writer http.ResponseWriter, r *http.Request) {
	config := configDto{DefaultSeed: defaultNoise.Seed}
	_ = json.NewEncoder(writer).Encode(config)
}
