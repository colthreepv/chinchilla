package main

import (
	"github.com/zenazn/goji"

	"log"
	"net/http"
	"os"
)

// serveStatic makes sure that if Chinchilla is running for development
// purposes, it serves static content from /static
// simple function with side-effects
func serveStatic() {
	buildMode := os.Getenv("chinchilla")
	switch buildMode {
	case "development", "":
		log.Println(info("Enabling static file serving dir: %s", "static/"))
		goji.Handle("/", http.FileServer(http.Dir("static")))
	case "production":
		return
	default:
		return
	}
}
