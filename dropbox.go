package main

import (
	"github.com/stacktic/dropbox"
	mgo "gopkg.in/mgo.v2"

	// "github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"encoding/json"
	"net/http"
)

type helloJson struct {
	DropboxUser string
}

// handler functions
func helloHandler(db *dropbox.Dropbox, s *mgo.Session) web.Handler {
	userC := s.DB(chi.Mongo.Database).C("User")
	gojiHandler := func(c web.C, w http.ResponseWriter, r *http.Request) {
		var h helloJson
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&h)
		if err != nil || h.DropboxUser == "" || len(h.DropboxUser) != 64 {
			http.Error(w, NewChiError("Invalid JSON received"), http.StatusBadRequest)
			return
		}
		// EVERYTHING WENT WELL, OH MY GOD AWESOME!
		err = userC.Insert(&h)
		if err != nil {
			http.Error(w, NewChiError(err.Error()), http.StatusBadRequest)
			return
		}
	}
	return web.HandlerFunc(gojiHandler)
}
