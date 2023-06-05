package main

import (
	"forum/internal/pkg/delivery"
	"forum/internal/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.JSONSerializer = &utils.EasyJSONSerializer{}

	db, err := utils.InitPostgres()
	if err != nil {
		e.Logger.Fatal(err)
	}

	api := e.Group("/api")
	api.Use(middleware.Recover(), middleware.Logger())

	users := delivery.NewUserHandler(db)
	service := delivery.NewServiceHandler(db)

	api.GET("/user/:nickname/profile", users.GetUser)
	api.POST("/user/:nickname/create", users.CreateUser)
	api.POST("/user/:nickname/profile", users.UpdateUser)

	api.GET("/service/status", service.Status)
	api.POST("/service/clear", service.Clear)

	e.Logger.Fatal(e.Start(":5000"))
}
