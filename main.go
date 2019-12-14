package main

import (
	"github.com/cyverse-de/echo-middleware/redoc"
	"github.com/cyverse-de/requests/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(redoc.Serve(redoc.Opts{Title: "DE Administrative Requests API Documentation"}))

	a := api.API{Echo: e}

	e.GET("/", a.RootHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
