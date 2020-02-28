package web

import (
	"fmt"
	"github.com/greboid/irc/config"
	"github.com/greboid/irc/database"
	"github.com/greboid/irc/irc"
	"github.com/labstack/echo"
	"github.com/tidwall/buntdb"
	"log"
	"net/http"
	"strings"
)

type Web struct {
	conf *config.Config
	irc  *irc.Connection
	db   *database.DB
}

func NewWeb(conf *config.Config, irc *irc.Connection, db *database.DB) *Web {
	log.Print("Initialising web")
	return &Web{
		conf: conf,
		irc:  irc,
		db:   db,
	}
}

func (web *Web) StartWeb() {
	log.Printf("Starting web: %d", web.conf.WebPort)
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	admin := e.Group("admin/:key", adminKeyMiddleware(web.conf))
	admin.GET("/keys", web.getKeys)
	key := e.Group("/user/:key", userKeyMiddleware(web.db))
	key.GET("/message/:message", web.sendMessage)
	err := e.Start(fmt.Sprintf("0.0.0.0:%d", web.conf.WebPort))
	if err != nil {
		log.Panicf("Unable to start web server: %v", err)
	}
	log.Print("Finished web")
}

func adminKeyMiddleware(conf *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Param("key") == conf.AdminKey {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}

func userKeyMiddleware(db *database.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if db.CheckKey(c.Param("key")) {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}

func (web *Web) sendMessage(context echo.Context) error {
	message := context.Param("message")
	go web.irc.SendRawf("PRIVMSG %s %s", web.conf.Channel, message)
	return context.String(http.StatusOK, "Done")
}

func (web *Web) getKeys(context echo.Context) error {
	var keys []string
	err := web.db.Db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			keys = append(keys, key)
			return true
		})
		return err
	})
	if err != nil {
		return context.String(http.StatusInternalServerError, "Error")
	}
	return context.String(http.StatusOK, strings.Join(keys, ", "))
}
