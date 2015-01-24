package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db := session.DB("test")
	index := mgo.Index{
		Key:    []string{"name"},
		Unique: true,
		Sparse: true,
	}
	if err := db.C("repos").EnsureIndex(index); err != nil {
		panic(err)
	}
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	repo := router.Group("/client/repo")
	{
		repo.POST("/:name", func(c *gin.Context) {
			repoName := c.Params.ByName("name")
			//TODO validate param
			newRepo := Repo{repoName, 0}
			err := db.C("repos").Insert(&newRepo)
			if err != nil {
				log.Println(err)
				c.String(500, "nok")
			} else {
				c.String(200, "OK")
			}
		})
		repo.GET("/", func(c *gin.Context) {
			result := []Repo{}
			err := db.C("repos").Find(bson.M{}).All(&result)
			if err != nil {
				log.Println(err)
				c.String(500, "nok")
			} else {
				c.JSON(200, result)
			}
		})
	}

	router.Run(":8080")
}
