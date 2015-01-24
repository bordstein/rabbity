package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
)

func main() {
	app := NewApp()
	if err := app.ConnectDatabase(); err != nil {
		panic(err)
	}
	defer app.DisconnectDatabase()

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
			err := app.DB.C("repos").Insert(&newRepo)
			if err != nil {
				log.Println(err)
				c.String(500, "nok")
			} else {
				c.String(200, "OK")
			}
		})
		repo.GET("/", func(c *gin.Context) {
			result := []Repo{}
			err := app.DB.C("repos").Find(bson.M{}).All(&result)
			if err != nil {
				log.Println(err)
				c.String(500, "nok")
			} else {
				c.JSON(200, result)
			}
		})
	}
	cluster := router.Group("/cluster")
	{
		cluster.GET("/fetch/:sha3sum", func(c *gin.Context) {
			sha3sum := c.Params.ByName("sha3sum")
			//TODO validate param
			file, err := app.Store.GetFile(sha3sum)
			if err != nil {
				log.Println(err)
				c.String(500, "nok")
			} else {
				io.Copy(c.Writer, file)
			}
		})
	}

	router.Run(":8080")
}
