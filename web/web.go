package web

import (
	"fmt"
	"github.com/greboid/irc/database"
	"github.com/greboid/irc/irc"
	"github.com/labstack/echo"
	"log"
)

type Web struct {
	irc      *irc.Connection
	db       *database.DB
	webPort  int
	channel  string
	adminKey string
}

func NewWeb(webPort int, channel string, adminKey string, irc *irc.Connection, db *database.DB) *Web {
	log.Print("Initialising web")
	return &Web{
		irc:      irc,
		db:       db,
		webPort:  webPort,
		channel:  channel,
		adminKey: adminKey,
	}
}

func (web *Web) StartWeb() {
	log.Printf("Starting web: %d", web.webPort)
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	addAdminRoutes(e, web)
	addBasicRoutes(e, web)
	err := e.Start(fmt.Sprintf("0.0.0.0:%d", web.webPort))
	if err != nil {
		log.Fatalf("Unable to start web server: %v", err)
	}
	log.Print("Finished web")
}
