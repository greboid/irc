package main

import (
	"github.com/labstack/echo"
)

func adminKeyMiddleware(adminKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Param("key") == adminKey {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}

func userKeyMiddleware(db *DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if db.CheckUser(c.Param("key")) {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}
