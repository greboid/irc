package database

import (
	"encoding/json"
	"github.com/tidwall/buntdb"
	"log"
)

type DB struct {
	path string
	Db   *buntdb.DB
}

type Plugin struct {
	Name  string
	Token string
}

type User struct {
	Name  string
	Token string
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

func (db *DB) setKey(key string, value string) error {
	err := db.Db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})
	return err
}

func (db *DB) getKey(key string) (string, error) {
	var value string
	var err error
	err = db.Db.View(func(tx *buntdb.Tx) error {
		dbValue, err := tx.Get(key)
		if err != nil {
			return err
		}
		value = dbValue
		return nil
	})
	return value, err
}

func (db *DB) getUsers() []User {
	dataString, err := db.getKey("users")
	if err != nil {
		return nil
	}
	var data []User
	err = json.Unmarshal([]byte(dataString), &data)
	if err != nil {
		return nil
	}
	return data
}

func (db *DB) getPlugins() []Plugin {
	dataString, err := db.getKey("plugins")
	if err != nil {
		return nil
	}
	var data []Plugin
	err = json.Unmarshal([]byte(dataString), &data)
	if err != nil {
		return nil
	}
	return data
}

func (db *DB) CheckUser(key string) bool {
	users := db.getUsers()
	for _, user := range users {
		if user.Token == key {
			return true
		}
	}
	return false
}

func (db *DB) CheckPlugin(key string) bool {
	plugins := db.getPlugins()
	for _, plugin := range plugins {
		if plugin.Token == key {
			return true
		}
	}
	return false
}

func (db *DB) CreateUser(name string, token string) error {
	user := User{
		Name:  name,
		Token: token,
	}
	users := append(db.getUsers(), user)
	value, err := json.Marshal(&users)
	if err != nil {
		return err
	}
	return db.setKey("users", string(value))
}

func (db *DB) CreatePlugin(name string, token string) error {
	plugin := Plugin{
		Name:  name,
		Token: token,
	}
	plugins := append(db.getPlugins(), plugin)
	value, err := json.Marshal(plugins)
	if err != nil {
		return err
	}
	return db.setKey("plugins", string(value))
}

func (db *DB) InitDB() {
	instance, err := buntdb.Open(db.path)
	if err != nil {
		log.Fatal(err)
	}
	db.Db = instance
}
