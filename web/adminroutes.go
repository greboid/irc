package web

import (
	"github.com/labstack/echo"
	"github.com/tidwall/buntdb"
	"net/http"
	"strings"
)

func addAdminRoutes(e *echo.Echo, web *Web) {
	admin := e.Group("admin/:key", adminKeyMiddleware(web.conf))
	admin.GET("/keys", web.getKeys)
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
