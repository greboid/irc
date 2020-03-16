package web

import (
	"context"
	"github.com/greboid/irc/bot"
	"github.com/greboid/irc/rpc"
	"github.com/labstack/echo"
	"net/http"
)

func addBasicRoutes(e *echo.Echo, web *Web) {
	user := e.Group("/user/:key", userKeyMiddleware(web.db))
	user.GET("/message/:message", web.sendMessage)
	user.POST("/github/webhook", web.github)
}

func (web *Web) sendMessage(ctx echo.Context) error {
	message := ctx.Param("message")
	_, _ = web.irc.SendChannelMessage(rpc.CtxWithToken(context.Background(), "bearer", web.token), &rpc.ChannelMessage{
		Channel: web.channel,
		Message: message,
	})
	return ctx.String(http.StatusOK, "Done")
}

func (web *Web) github(context echo.Context) (err error) {
	wh := new(bot.GitHubCommitHook)
	if err = context.Bind(wh); err != nil {
		return
	}
	return context.JSON(http.StatusOK, wh)
}
