package main

import (
	"net/http"

	"forum/internal/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.JSONSerializer = &utils.EasyJSONSerializer{}

	api := e.Group("/api")
	api.Use(middleware.Recover(), middleware.Logger())

	api.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":5000"))
}
