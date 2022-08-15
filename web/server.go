package web

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/renderer"
	"github.com/fafeitsch/go-infinite-rail-generator/world"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type ApiOptions struct {
	Shift        int
	Seed         string
	TownNames    []string
	BuildVersion string
	BuildTime    string
}

func ApiHandler(options ApiOptions) http.HandlerFunc {
	generator := world.NewGenerator(options.Seed)
	generator.TownNames = options.TownNames
	return func(writer http.ResponseWriter, r *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		var head string
		originalPath := r.URL.Path
		head, r.URL.Path = shiftPath(r.URL.Path)
		switch head {
		case "tiles":
			serveTile(generator, options.Shift, writer, r)
			break
		case "config":
			serveConfig(options, writer, r)
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

func serveTile(defaultGenerator *world.Generator, shift int, writer http.ResponseWriter, r *http.Request) {
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
	var aNoise *world.Generator
	if seedString == "" || seedString == defaultGenerator.Seed() {
		aNoise = defaultGenerator
	} else {
		aNoise = world.NewGenerator(seedString)
		aNoise.TownNames = defaultGenerator.TownNames
	}
	tile := aNoise.Generate(hectometer + shift)
	writer.Header().Set("Content-Type", "image/svg+xml")
	rn := renderer.New(writer, 200)
	err = rn.Render(tile)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

type configDto struct {
	Version     string `json:"version"`
	BuildTime   string `json:"buildTime"`
	DefaultSeed string `json:"defaultSeed"`
}

func serveConfig(options ApiOptions, writer http.ResponseWriter, r *http.Request) {
	config := configDto{
		DefaultSeed: options.Seed,
		BuildTime:   options.BuildTime,
		Version:     options.BuildVersion,
	}
	_ = json.NewEncoder(writer).Encode(config)
}
