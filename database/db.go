package database

import (
	"github.com/tidwall/buntdb"
	"log"
)

type DB struct {
	path string
	Db   *buntdb.DB
}

func New(path string) *DB {
	log.Print("Initialising database")
	db := &DB{
		path: path,
		Db:   nil,
	}
	db.InitDB()
	var config buntdb.Config
	if err := db.Db.ReadConfig(&config); err != nil {
		log.Fatal(err)
	}
	config.SyncPolicy = buntdb.Always
	if err := db.Db.SetConfig(config); err != nil {
		log.Fatal(err)
	}
	log.Print("Initialised database")
	return db
}

func (db *DB) CheckKey(key string) bool {
	var valid bool
	err := db.Db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(key)
		if err != nil {
			return err
		}
		valid = true
		return nil
	})
	if err != nil {
		return false
	}
	return valid
}

func (db *DB) CreateKey(key string) error {
	log.Printf("Creating key: %s", key)
	err := db.Db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, key, nil)
		return err
	})
	return err
}

func (db *DB) InitDB() {
	instance, err := buntdb.Open(db.path)
	if err != nil {
		log.Fatal(err)
	}
	db.Db = instance
}
