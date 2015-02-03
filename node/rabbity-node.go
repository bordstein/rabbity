package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"net/http"
)

func main() {
	app := NewApp()
	defer app.DisconnectDatabase()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.POST("/client/upload", func(c *gin.Context) {
		reader, err := c.Request.MultipartReader()
		hashsum := ""
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		for true {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			switch part.FormName() {
			case "payload":
				hashsum, err = app.Store.AddFile(part)
				if err != nil {
					c.JSON(http.StatusInternalServerError,
						err.Error())
					return
				}
				break
			default:
				log.Println("ignoring unknown part: ",
					part.FormName())
			}
		}

		if hashsum != "" {
			c.JSON(http.StatusOK, hashsum)
		} else {
			c.JSON(http.StatusBadRequest, "file part missing")
		}

	})
	repo := router.Group("/client/repo")
	{
		repo.POST("/:name", func(c *gin.Context) {
			repoName := c.Params.ByName("name")
			//TODO validate param
			newRepo := Repo{repoName, 0}
			repoCol, err := app.GetRepoCol()
			if err != nil {
				c.String(500, err.Error())
				return
			}

			err = repoCol.Insert(&newRepo)
			if err != nil {
				log.Println(err)
				c.String(500, "nok")
			} else {
				c.String(200, "OK")
			}
		})
		repo.GET("/", func(c *gin.Context) {
			result := []Repo{}
			repoCol, err := app.GetRepoCol()
			if err != nil {
				c.String(500, err.Error())
				return
			}
			err = repoCol.Find(bson.M{}).All(&result)
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
		cluster.POST("/init", func(c *gin.Context) {
			initMsg := NodeInitMsg{}
			ok := c.BindWith(&initMsg, binding.JSON)
			if !ok {
				c.String(500, "could not parse body")
				return
			}
			err := initiateCluster(initMsg.Host)
			if err != nil {
				c.String(500, err.Error())
				return
			}
			c.String(200, "ok")
		})
	}

	router.Run(":8080")
}
