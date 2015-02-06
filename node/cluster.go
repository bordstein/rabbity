package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func initiateCluster(host string) error {
	cfg := bson.M{
		"_id": "rs0",
		"members": []bson.M{
			bson.M{"_id": 0, "host": host, "tags": bson.M{"rabbity_port": "8080"}},
		},
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

func getNodeInfo() (interface{}, error) {
	dbSession, err := mgo.DialWithInfo(
		&mgo.DialInfo{
			Addrs:    []string{"localhost"},
			Direct:   true,
			FailFast: true,
			Database: "local",
		})
	if err != nil {
		return nil, err
	}

	var result []bson.M
	err = dbSession.DB("local").C("system.replset").Find(bson.M{}).All(&result)
	log.Println(result)
	return result[0]["members"], err
}
