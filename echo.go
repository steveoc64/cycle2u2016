package main

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var e *echo.Echo

func _initEcho() {
	e = echo.New()
	e.Index("./cycle2u.com.au/index.html")
	e.ServeDir("/", "./cycle2u.com.au")

	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(mw.Gzip())
	if Config.Debug {
		e.SetDebug(true)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowCredentials: true,
		Debug:            Config.Debug,
	})
	e.Use(c.Handler)
}
