package main

import (
	mgo "gopkg.in/mgo.v2"

	"encoding/json"
	"log"
	"time"
)

type ChiError struct {
	Message string
}

func NewChiError(msg string) string {
	jsonMessage, _ := json.Marshal(ChiError{Message: msg})
	return string(jsonMessage)
}

type EngineError struct {
	E         error
	CreatedAt time.Time
	User      string // HARD: make a relationship into mongo, I DONT KNOW YET!!!
}

// FIXME: errors does get populated correctly, `e` is always empyy
func NewMongoError(u *ChiUser, mdb *mgo.Database, e error) {
	errorC := mdb.C("Error")
	var mErr *EngineError = &EngineError{E: e, CreatedAt: time.Now(), User: u.DropboxUser}
	err := errorC.Insert(mErr)
	if err != nil { // YO DAWN AN ERROR IN AN ERROR HANDLER.. WTF?
		log.Fatalf(debug(`Error occurred while trying to post an EngineError to mongo\n
EngineError: %#v\n
Error: %#v\n`, mErr, err))
	}
}
