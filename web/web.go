package web

import (
	"fmt"
	"github.com/greboid/irc/config"
	"github.com/greboid/irc/database"
	"github.com/greboid/irc/irc"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

type Web struct {
	conf *config.Config
	irc  *irc.Connection
	db   *database.DB
}

func NewWeb(conf *config.Config, irc *irc.Connection, db *database.DB) *Web {
	return &Web{
		conf: conf,
		irc:  irc,
		db:   db,
	}
}

func (web *Web) StartWeb() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.GET("/message/:key", web.testMessage)
	err := e.Start(fmt.Sprintf("0.0.0.0:%d", web.conf.WebPort))
	if err != nil {
		log.Panicf("Unable to start web server: %v", err)
	}
}

func (web *Web) testMessage(context echo.Context) error {
	key := context.Param("key")
	if web.db.CheckKey(key) {
		go web.irc.SendRawf("PRIVMSG %s beep", web.conf.Channel)
		return context.String(http.StatusOK, "Done")
	} else {
		return context.String(http.StatusForbidden, "Fordibben")
	}
}
