package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"os"
	//"log"
)

func main1() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	gridFs := session.DB("test").GridFS("fs")

	query := gridFs.Find(nil).Sort("filename")
	iter := query.Iter()
	var f *mgo.GridFile
	for gridFs.OpenNext(iter, &f) {
		fmt.Printf("Filename: %s\n", f.Name())
	}
	if iter.Close() != nil {
		panic(iter.Close())
	}

	var gFile *bson.D
	err = gridFs.Find(nil).One(&gFile)
	check(err)
	fmt.Println(gFile.Map()["filename"])
	fmt.Println(gFile)
	f, err = gridFs.OpenId(gFile.Map()["_id"])
	check(err)
	fmt.Println(f.Name())
	fmt.Println(f.MD5())
	fmt.Println("here comes the content:\n=====")
	io.Copy(os.Stdout, f)
	fmt.Println("\n=====")
}

func check1(err error) {
	if err != nil {
		panic(err)
	}
}
