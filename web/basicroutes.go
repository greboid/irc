package web

import (
	"github.com/greboid/irc/bot"
	"github.com/labstack/echo"
	"net/http"
)

func addBasicRoutes(e *echo.Echo, web *Web) {
	user := e.Group("/user/:key", userKeyMiddleware(web.db))
	user.GET("/message/:message", web.sendMessage)
	user.POST("/github/webhook", web.github)
}

func (web *Web) sendMessage(context echo.Context) error {
	message := context.Param("message")
	go web.irc.SendRawf("PRIVMSG %s :%s", web.channel, message)
	return context.String(http.StatusOK, "Done")
}

func (web *Web) github(context echo.Context) (err error) {
	wh := new(bot.GitHubCommitHook)
	if err = context.Bind(wh); err != nil {
		return
	}
	return context.JSON(http.StatusOK, wh)
}
