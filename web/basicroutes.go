package web

import (
	"github.com/labstack/echo"
	"net/http"
)

func addBasicRoutes(e *echo.Echo, web *Web) {
	user := e.Group("/user/:key", userKeyMiddleware(web.db))
	user.GET("/message/:message", web.sendMessage)
}

func (web *Web) sendMessage(context echo.Context) error {
	message := context.Param("message")
	go web.irc.SendRawf("PRIVMSG %s %s", web.conf.Channel, message)
	return context.String(http.StatusOK, "Done")
}