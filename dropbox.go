package main

import (
	"github.com/stacktic/dropbox"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	// "github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ChiUser struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	DropboxUser  string        `bson:"dropbox_user"`
	LatestCursor string        `bson:"latest_cursor"`
}

func (u *ChiUser) UpdateCursor(mdb *mgo.Database, cursor string) {
	update := bson.M{"$set": bson.M{"latest_cursor": cursor}}
	err := mdb.C("User").UpdateId(u.ID, update)
	if err != nil {
		NewMongoError(u, mdb, err)
	}
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

		// TODO: use Upsert

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
		h.ID = bson.NewObjectId()
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

type ChiImage struct {
	User  string
	Path  string
	Entry *dropbox.Entry
}

// Struct to track how much a specific Operation takes
type ChiStat struct {
	StartTime      time.Time     `bson:"start_time"`
	ElapsedTime    time.Duration `bson:"elapsed_time"`
	Operation      string
	OperationCount int `bson:"operation_count"`
}

type Downloader struct {
	u   *ChiUser
	db  *dropbox.Dropbox
	mdb *mgo.Database
	// something to track goroutines created???
}

// constructor for Downloader
func NewDownloader(u *ChiUser, db *dropbox.Dropbox, mdb *mgo.Database) *Downloader {
	return &Downloader{u: u, db: db, mdb: mdb}
}

func (d Downloader) Start() {
	defer d.Continue("") // we call Continue starting with a blank cursor
	log.Println(debug("Starting image loading"))
}

func (d Downloader) Continue(cursor string) {
	var stat *ChiStat = &ChiStat{StartTime: time.Now(), Operation: "Delta"}
	d.db.SetAccessToken(d.u.DropboxUser)
	// a delta call with an empty cursor, as described here
	// https://www.dropbox.com/developers/blog/69/efficiently-enumerating-dropbox-with-delta
	deltaP, err := d.db.Delta(cursor, "/")
	if err != nil {
		NewMongoError(d.u, d.mdb, err)
		return // ends this goroutine with extreme failure
	}
	if len(deltaP.Entries) >= 0 { // add all images reported from dropbox, to mongo
		images := d.mdb.C("Image").Bulk()
		for _, dEntry := range deltaP.Entries {
			images.Insert(&ChiImage{User: d.u.DropboxUser, Path: dEntry.Path, Entry: dEntry.Entry})
		}
		images.Unordered()
		_, err := images.Run()
		if err != nil {
			NewMongoError(d.u, d.mdb, err)
			return
		}
	}
	if deltaP.HasMore {
		defer d.Continue(deltaP.Cursor.Cursor)
	} else {
		// update user with latest cursor
		d.u.UpdateCursor(d.mdb, deltaP.Cursor.Cursor)
		// TODO: notify that delta is done to the ThumbnailDownloader
		log.Println(debug("Wooo! All files bulk-inserted!"))
	}

	// Send statistical data
	stat.ElapsedTime = time.Since(stat.StartTime)
	stat.OperationCount = len(deltaP.Entries)
	err = d.mdb.C("Stat").Insert(stat)
	if err != nil {
		NewMongoError(d.u, d.mdb, err)
		return
	}
}

func downloaderRoutine(u chan *ChiUser, db *dropbox.Dropbox, s *mgo.Session) {
	for {
		var newUser *ChiUser
		newUser = <-u
		log.Println(debug("routine notified: %+v", newUser))
		newDownloader := NewDownloader(newUser, db, s.DB(chi.Mongo.Database))
		go newDownloader.Start() // on user creation a goroutine gets assigned to a user (quick n dirty)
	}
}
