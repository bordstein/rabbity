package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Phone string
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db := session.DB("test")
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.GET("/new", func(c *gin.Context) {
		err := db.C("people").Insert(&Person{"Ale", "+55 53 8116 9639"},
			&Person{"Cla", "+55 53 8402 8510"})
		if err != nil {
			fmt.Println(err)
			c.String(500, "nok")
		} else {
			c.String(200, "OK")
		}
	})
	router.GET("/ls", func(c *gin.Context) {
		result := Person{}
		err := db.C("people").Find(bson.M{"name": "Ale"}).One(&result)
		if err != nil {
			fmt.Println(err)
			c.String(500, "nok")
		} else {
			c.JSON(200, result)
		}
	})

	router.Run(":8080")
}
