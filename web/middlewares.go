package web

import (
	"github.com/greboid/irc/config"
	"github.com/greboid/irc/database"
	"github.com/labstack/echo"
)

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
			if db.CheckUser(c.Param("key")) {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}
