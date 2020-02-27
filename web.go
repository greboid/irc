package main

import (
	"fmt"
	"github.com/greboid/irc/irc"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

type Web struct {
	channel string
	port    int
	irc     *irc.IRCConnection
}

func (web *Web) StartWeb() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.GET("/", web.rootPath)
	err := e.Start(fmt.Sprintf("0.0.0.0:%d", web.port))
	if err != nil {
		log.Panicf("Unable to start web server: %v", err)
	}
}

func (web *Web) rootPath(context echo.Context) error {
	go web.irc.SendRawf("PRIVMSG %s beep", web.channel)
	return context.String(http.StatusOK, "Done")
}
