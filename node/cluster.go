package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func initiateCluster(host string) error {
	cfg := bson.D{
		{"_id", "rs0"},
		{"members", []bson.D{
			bson.D{{"_id", 0}, {"host", host}},
		}},
	}

	dbSession, err := mgo.DialWithInfo(
		&mgo.DialInfo{
			Addrs:    []string{"localhost"},
			Direct:   true,
			FailFast: true,
			Database: "admin",
		})
	if err != nil {
		return err
	}
	dbSession.SetMode(mgo.Monotonic, true)

	result := bson.D{}
	log.Print(bson.D{{"replSetInitiate", cfg}})
	err = dbSession.Run(bson.D{{"replSetInitiate", cfg}}, &result)
	return err
}
