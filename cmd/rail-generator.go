package main

import (
	"github.com/fafeitsch/go-infinite-rail-generator/web"
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe("0.0.0.0:9551", web.ApiHandler("go-infinite-rail-generator"))
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
