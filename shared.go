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

func NewMongoError(u *ChiUser, s *mgo.Session, e error) {
	errorC := s.DB(chi.Mongo.Database).C("Error")
	var mErr *EngineError = &EngineError{E: e, CreatedAt: time.Now(), User: u.DropboxUser}
	err := errorC.Insert(mErr)
	if err != nil { // YO DAWN AN ERROR IN AN ERROR HANDLER.. WTF?
		log.Fatalf(debug(`Error occurred while trying to post an EngineError to mongo\n
EngineError: %#v\n
Error: %#v\n`, mErr, err))
	}
}
