package web

import (
	"fmt"
	"github.com/greboid/irc/config"
	"github.com/greboid/irc/database"
	"github.com/greboid/irc/irc"
	"github.com/labstack/echo"
	"log"
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
	addAdminRoutes(e, web)
	addBasicRoutes(e, web)
	err := e.Start(fmt.Sprintf("0.0.0.0:%d", web.conf.WebPort))
	if err != nil {
		log.Panicf("Unable to start web server: %v", err)
	}
	log.Print("Finished web")
}
