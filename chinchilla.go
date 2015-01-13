package main

import (
	"github.com/stacktic/dropbox"
	mgo "gopkg.in/mgo.v2"

	_ "github.com/pressly/cji"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"log"
)

var chi *chiConfig
var chiMongo *mgo.Session

func main() {
	chi = newChiConfig("config.toml")
	db := dropbox.NewDropbox()
	db.SetAppInfo(chi.Dropbox.Key, chi.Dropbox.Secret)
	chiMongo, err := mgo.DialWithInfo(&chi.Mongo)
	if err != nil {
		panic(err)
	}
	// databases, _ := chiMongo.DatabaseNames()
	// log.Printf("chiMongo: %+v", databases)
	log.Println(db)

	// start goji
	api := web.New()
	goji.Handle("/api/*", api)

	// middlewares in use for any /api/ route
	api.Use(middleware.SubRouter)
	api.Use(middleware.EnvInit)
	// api.Use(headerCheck)

	// specific /api/:name path
	// api.Get("/test/:name", cji.Use(fakeDatabaseReq).On(printName))
	// api.Get("/search", listImages(db))
	api.Post("/hello", helloHandler(db, chiMongo))

	goji.Serve()
}
