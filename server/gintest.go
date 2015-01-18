package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"io"
	"log"
	"nurio.at/testi"
	"os"
)

type Person struct {
	Name string
	Age  int
}

type App struct {
	DB *mgo.Database
}

func (app App) HandleFileUpload(c *gin.Context) {
	gridFs := app.DB.GridFS("fs")
	file := new(mgo.GridFile)
	metadata := new(testi.FileMetaData)

	reader, err := c.Request.MultipartReader()
	if err != nil {
		panic(err)
	}
	for true {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		log.Println("loop..")
		switch part.FormName() {
		case "payload":
			file, err = gridFs.Create(part.FileName())
			if err != nil {
				panic(err)
			}
			_, sha3sum, err := testi.Sha3HashCopy(file, part)
			if err != nil {
				panic(err)
			}

			metadata.Name = part.FileName()
			metadata.Sha3sum = sha3sum
			break
		default:
			log.Println("ignoring unknown part: ",
				part.FormName())
		}
	}

	file.SetMeta(metadata)
	err = file.Close()
	if err != nil {
		panic(err)
	}
	c.JSON(200, metadata)
}

func main() {
	app := new(App)
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	app.DB = session.DB("test")
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		p := new(Person)
		p.Age = 5
		p.Name = "hans"
		c.JSON(200, p)
	})
	router.GET("/filetest", func(c *gin.Context) {
		file, err := os.Open("test.go") // For read access.
		if err != nil {
			c.Fail(500, err)
		}
		io.Copy(c.Writer, file)
	})
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.POST("/upload", app.HandleFileUpload)
	router.POST("/submit", func(c *gin.Context) {
		c.String(401, "not authorized")
	})
	router.PUT("/error", func(c *gin.Context) {
		c.String(500, "and error hapenned :(")
	})
	router.Run(":8080")
}
