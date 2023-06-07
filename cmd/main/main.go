package main

import (
	"forum/internal/pkg/delivery"
	"forum/internal/pkg/usecase"
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
	api.Use(middleware.Recover())

	threadsUsecase := usecase.NewThreadUsecase(db)
	forumsUsecase := usecase.NewForumUsecase(db)

	users := delivery.NewUserHandler(db, forumsUsecase)
	service := delivery.NewServiceHandler(db)
	forums := delivery.NewForumHandler(forumsUsecase)
	threads := delivery.NewThreadHandler(threadsUsecase)
	posts := delivery.NewPostHandler(db, threadsUsecase, forumsUsecase)
	votes := delivery.NewVoteHandler(db, threadsUsecase)

	api.GET("/user/:nickname/profile", users.GetUser)
	api.POST("/user/:nickname/create", users.CreateUser)
	api.POST("/user/:nickname/profile", users.UpdateUser)
	api.GET("/forum/:slug/users", users.GetByForum)

	api.POST("/forum/create", forums.Create)
	api.GET("/forum/:slug/details", forums.Get)

	api.POST("/forum/:slug/create", threads.Create)
	api.GET("/forum/:slug/threads", threads.GetByForum)
	api.GET("/thread/:slug_or_id/details", threads.Get)
	api.POST("/thread/:slug_or_id/details", threads.Update)

	api.POST("/thread/:slug_or_id/create", posts.AddPosts)
	api.GET("/post/:id/details", posts.GetPost)
	api.POST("/post/:id/details", posts.UpdatePost)
	api.GET("/thread/:slug_or_id/posts", posts.GetPosts)

	api.POST("/thread/:slug_or_id/vote", votes.AddVote)

	api.GET("/service/status", service.Status)
	api.POST("/service/clear", service.Clear)

	e.Logger.Fatal(e.Start(":5000"))
}
