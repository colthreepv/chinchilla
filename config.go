package main

import (
	"github.com/BurntSushi/toml"
	mgo "gopkg.in/mgo.v2"
)

type chiConfig struct {
	Dropbox dropboxConfig
	Mongo   mgo.DialInfo
}
type dropboxConfig struct {
	Key, Secret string
}

func newChiConfig(filePath string) *chiConfig {
	var c chiConfig
	if meta, err := toml.DecodeFile(filePath, &c); err == nil {
		if meta.IsDefined("dropbox", "key") &&
			meta.IsDefined("dropbox", "secret") &&
			meta.IsDefined("mongo", "addrs") &&
			meta.IsDefined("mongo", "database") {
			return &c
		} else {
			panic("the toml file must provide enough informations to fill dropboxConfig and mgo.DialInfo")
		}
		panic(err)
	}
	panic("a valid config.toml is required to start")
}
