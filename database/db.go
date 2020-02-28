package database

import (
	"github.com/tidwall/buntdb"
	"log"
)

type DB struct {
	path string
	db   *buntdb.DB
}

func New(path string) *DB {
	db := &DB{
		path: path,
		db:   nil,
	}
	db.InitDB()
	var config buntdb.Config
	if err := db.db.ReadConfig(&config); err != nil {
		log.Fatal(err)
	}
	config.SyncPolicy = buntdb.Always
	if err := db.db.SetConfig(config); err != nil {
		log.Fatal(err)
	}
	return db
}

func (db *DB) CheckKey(key string) bool {
	var valid bool
	err := db.db.View(func(tx *buntdb.Tx) error {
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
	err := db.db.Update(func(tx *buntdb.Tx) error {
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
	db.db = instance
}
