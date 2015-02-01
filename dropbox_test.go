package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"log"
	"testing"
)

func TestChiUser(t *testing.T) {
	// tear-up code (copied from main)
	chi := newChiConfig("config.toml")
	log.Println(info("Trying to establish connection with mongoDB server %s", chi.Mongo.Addrs))
	chiMongo, err := mgo.DialWithInfo(&chi.Mongo)
	if err != nil {
		t.Error(err)
	}
	defer chiMongo.Close()
	mdb := chiMongo.DB(chi.Mongo.Database)

	// test code
	u := &ChiUser{ID: bson.NewObjectId(), DropboxUser: "SOME-ID", mongoCollection: mdb.C("User")}
	err = u.Save()
	if err != nil {
		t.Error(err)
	}
	err = u.UpdateCursor("SOME-CURSOR")
	if err != nil {
		t.Error(err)
	}

	// tear-down code
	mdb.C("User").RemoveId(u.ID)
}
