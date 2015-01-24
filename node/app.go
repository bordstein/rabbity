package main

import (
	"gopkg.in/mgo.v2"
)

type App struct {
	DB        *mgo.Database
	Store     FStore
	dbSession *mgo.Session
}

func NewApp() App {
	store := FStore{Path: "/tmp/test", TmpPath: "/tmp/test/tmp"}
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
