package main

import (
	"github.com/BurntSushi/toml"
	// "github.com/codegangsta/negroni"
	// "github.com/julienschmidt/httprouter"
	"github.com/stacktic/dropbox"

	"github.com/pressly/cji"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"fmt"
	"log"
	"net/http"
	"time"
)

type config struct {
	DropboxAppKey, DropboxAppSecret string
}

func NewConfig(filePath string) *config {
	var c config
	if meta, err := toml.DecodeFile(filePath, &c); err == nil {
		if meta.IsDefined("DropboxAppKey") && meta.IsDefined("DropboxAppSecret") {
			return &c
		} else {
			panic("the toml file must provide DropboxAppKey and DropboxAppSecret keys")
		}
		panic(err)
	}
	panic("a valid config.toml is required to start")
}

func headerDumpMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if customHeader := r.Header.Get("X-CUSTOM"); customHeader != "" {
		log.Printf("received request with custom Header: %s", customHeader)
	}
	next(rw, r)
}

func headerCheck(c *web.C, h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if dropboxUser := r.Header.Get("X-Dropbox-User"); dropboxUser != "" && len(dropboxUser) == 64 {
			// header syntactically correct
			log.Printf("dropbox user: %s", dropboxUser)
			c.Env["dropboxUser"] = dropboxUser
			h.ServeHTTP(w, r)
		} else {
			http.Error(w, "missing X-Dropbox-User Header", 400)
		}
	}
	return http.HandlerFunc(handler)
}

func responseGivingMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Header().Set("Content-Type", "json")
	rw.Write([]byte(`{ "message": "all is fine!" }`))
}

// trying a chain of 2 middlewares
func fakeDatabaseReq(c *web.C, h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * 100 * time.Millisecond)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}

func printName(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Happy to see you %s, your dropboxUser token is: %s", c.URLParams["name"], c.Env["dropboxUser"])
}

func main() {
	// standard init dropbox package for pics-or-stfu
	db := dropbox.NewDropbox()
	c := NewConfig("config.toml")
	db.SetAppInfo(c.DropboxAppKey, c.DropboxAppSecret)

	/*
		// router is an http.Handler - compatible structure
		router := httprouter.New()
		// but also negroni.New is an http.Handler - compatible structure!
		awesomeMiddlewareStack := negroni.New(negroni.HandlerFunc(headerDumpMiddleware), negroni.HandlerFunc(responseGivingMiddleware))
		router.Handler("GET", "/", awesomeMiddlewareStack)

		router.Handle("GET", "/:something", func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			fmt.Fprintf(rw, "hello, %s!\n", ps.ByName("something"))
		})

		classic := negroni.Classic()
		classic.UseHandler(router)
		classic.Run(":8080")
	*/

	api := web.New()
	goji.Handle("/api/*", api)

	// middlewares in use for any /api/ route
	api.Use(middleware.SubRouter)
	api.Use(middleware.EnvInit)
	api.Use(headerCheck)

	// specific /api/:name path
	api.Get("/:name", cji.Use(fakeDatabaseReq).On(printName))

	goji.Serve()

	/*
			db.SetAccessToken("****************************************************************")

			startTime := time.Now()
			images, err := db.Search("", ".jpg", 1000, true)
			dropboxElapsed := time.Since(startTime)
			if err != nil {
				log.Fatalln(err)
			}
			var thumbAvailable int
			for _, img := range images {
				if img.ThumbExists {
					thumbAvailable = thumbAvailable + 1
				}
			}
			log.Printf(
				`images with thumbnail available: %d
		    time elapsed: %s`, thumbAvailable, dropboxElapsed.String())
	*/
}
