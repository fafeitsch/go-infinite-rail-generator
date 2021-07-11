package main

import (
	"github.com/fafeitsch/go-infinite-rail-generator/noise"
	"github.com/fafeitsch/go-infinite-rail-generator/web"
	"log"
	"net/http"
)

func main() {
	defaultNoise := noise.New("Rails are fun!")
	err := http.ListenAndServe("0.0.0.0:9551", web.ApiHandler(defaultNoise))
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
