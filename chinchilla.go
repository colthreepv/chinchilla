package main

import (
	"github.com/stacktic/dropbox"
	mgo "gopkg.in/mgo.v2"

	// "github.com/pressly/cji"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/fatih/color"

	"log"
)

// colorized logger
var info = color.New(color.FgYellow).SprintfFunc()
var success = color.New(color.FgGreen).SprintfFunc()
var debug = color.New(color.FgRed).SprintfFunc()

// ChiRuntime is global struct meant to keep some
// mongoDB informations available everywhere
type ChiRuntime struct {
	Mongo           *mgo.Session
	ErrorCollection *mgo.Collection
}

var config *chiConfig
var Chi ChiRuntime

func main() {
	// importing configuration from .toml file
	// HARD dependencies of chinchilla are: dropbox app keys, mongodb
	config = newChiConfig("config.toml")
	db := dropbox.NewDropbox()
	db.SetAppInfo(config.Dropbox.Key, config.Dropbox.Secret)
	log.Println(info("Trying to establish connection with mongoDB server %s", config.Mongo.Addrs))
	mongoSession, err := mgo.DialWithInfo(&config.Mongo)
	if err != nil {
		panic(err)
	}
	Chi.Mongo = mongoSession
	defer mongoSession.Close()
	log.Println(success("mongoDB success"))

	var newUsers = make(chan *ChiUser)
	go downloaderRoutine(newUsers, db, Chi.Mongo)

	// start goji
	serveStatic()
	api := web.New()
	goji.Handle("/api/*", api)

	// middlewares in use for any /api/ route
	api.Use(middleware.SubRouter)
	api.Use(middleware.EnvInit)
	// api.Use(headerCheck)

	// specific /api/:name path
	// api.Get("/test/:name", cji.Use(fakeDatabaseReq).On(printName))
	// api.Get("/search", listImages(db))
	api.Post("/hello", helloHandler(db, Chi.Mongo, newUsers))

	goji.Serve()
}
