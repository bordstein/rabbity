package main

import (
	"flag"
	"gopkg.in/mgo.v2"
	"log"
	"path/filepath"
)

type App struct {
	DB        *mgo.Database
	Store     FStore
	dbSession *mgo.Session
}

func NewApp() App {
	dataDir := flag.String("datadir", "/tmp/rabbity-store",
		"The directory to store the binary files")
	flag.Parse()
	store := FStore{
		Path:    *dataDir,
		TmpPath: filepath.Join(*dataDir, "tmp"),
	}
	log.Printf("using datastore in %v", store.Path)
	return App{Store: store}
}

func (app *App) ConnectDatabase() error {
	var err error
	app.dbSession, err = mgo.Dial("localhost")
	if err != nil {
		return err
	}
	app.DB = app.dbSession.DB("test")
	index := mgo.Index{
		Key:    []string{"name"},
		Unique: true,
		Sparse: true,
	}
	if err := app.DB.C("repos").EnsureIndex(index); err != nil {
		return err
	}
	return nil
}

func (app *App) DisconnectDatabase() {
	app.dbSession.Close()
}
