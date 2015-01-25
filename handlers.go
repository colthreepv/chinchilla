package main

import (
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"log"
	"net/http"
	"os"
	"strings"
)

// ReverseStringPattern is a simple struct keeping the raw string
// and the prefix to skip in requests
type ReverseStringPattern struct {
	raw    string
	prefix string
}

// prefix returns "" so goji doesn't ever skip the matching function
func (s ReverseStringPattern) Prefix() string {
	return ""
}
func (s ReverseStringPattern) Match(r *http.Request, c *web.C) bool {
	path := r.URL.Path
	return !strings.HasPrefix(path, s.prefix)
}
func (s ReverseStringPattern) Run(r *http.Request, c *web.C) {}

func NewReverseStringPattern(s string) ReverseStringPattern {
	var prefix string = s
	if strings.HasSuffix(s, "*") {
		prefix = s[:len(s)-1]
	}
	return ReverseStringPattern{raw: s, prefix: prefix}
}

// serveStatic makes sure that if running for development
// purposes, it serves static content from /static
func serveStatic() {
	buildMode := os.Getenv("chinchilla")
	switch buildMode {
	case "development", "":
		log.Println(info("Enabling static file serving dir: %s", "static/"))
		reverseApiPattern := NewReverseStringPattern("/api/*")
		goji.Handle(reverseApiPattern, http.FileServer(http.Dir("static")))
	case "production":
		return
	default:
		return
	}
}
