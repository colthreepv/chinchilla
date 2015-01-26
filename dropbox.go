package main

import (
	"github.com/stacktic/dropbox"
	mgo "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"

	// "github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"encoding/json"
	"log"
	"net/http"
)

type ChiUser struct {
	DropboxUser string
}

// handler functions
func helloHandler(db *dropbox.Dropbox, s *mgo.Session, notify chan *ChiUser) web.Handler {
	userC := s.DB(chi.Mongo.Database).C("User")
	gojiHandler := func(c web.C, w http.ResponseWriter, r *http.Request) {
		var h ChiUser
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&h)
		if err != nil || h.DropboxUser == "" || len(h.DropboxUser) != 64 {
			http.Error(w, NewChiError("Invalid JSON received"), http.StatusBadRequest)
			return
		}

		// check if corrisponding user exists already
		isPresent, err := userC.Find(&h).Count()
		if err != nil {
			http.Error(w, NewChiError(err.Error()), http.StatusBadRequest)
			return
		}
		if isPresent > 0 { // means the user is already registered
			w.WriteHeader(http.StatusOK)
			return
		}
		// create user instead
		err = userC.Insert(&h)
		if err != nil {
			http.Error(w, NewChiError(err.Error()), http.StatusBadRequest)
			return
		}
		notify <- &h
		w.WriteHeader(http.StatusCreated)
	}
	return web.HandlerFunc(gojiHandler)
}

type Downloader struct {
	db *dropbox.Dropbox
	s  *mgo.Session
	// something to track goroutines created???
}

func fromDropToMongo(u *ChiUser, db *dropbox.Dropbox, s *mgo.Session) {
	db.SetAccessToken(u.DropboxUser)
	// a delta call with an empty cursor, as described here
	// https://www.dropbox.com/developers/blog/69/efficiently-enumerating-dropbox-with-delta
	root, err := db.Delta("", "/")
	if err != nil {
		NewMongoError(u, s, err)
		return // ends this goroutine with extreme failure
	}
	if len(root.Entries) >= 0 { // there's something to add to mongo
		// BULK INSERT HERE ----
	}
	log.Println(debug("delta : %#v", root))
	// TODO: this should not end here...
}

func downloaderRoutine(u chan *ChiUser, db *dropbox.Dropbox, s *mgo.Session) {
	for {
		var newUser *ChiUser
		newUser = <-u
		log.Println(debug("routine notified: %+v", newUser))
		go fromDropToMongo(newUser, db, chiMongo) // on user creation a goroutine gets assigned to a user (quick n dirty)
	}
}
