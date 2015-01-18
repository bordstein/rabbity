package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	//"log"
)

func _main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	gridFs := session.DB("test").GridFS("fs")

	file, err := gridFs.Create("myfile.txt")
	check(err)
	n, err := file.Write([]byte("Hello world!"))
	check(err)
	err = file.Close()
	check(err)
	fmt.Printf("%d bytes written\n", n)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
