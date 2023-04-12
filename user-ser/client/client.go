package main

import (
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/web"
	"user-ser/client/handler"
)

func main() {
	router := gin.Default()

	router.Handle("GET", "toLogin", func(context *gin.Context) {
		context.String(200, "to login...")
	})

	router.GET("/login", handler.UserLogin)

	service := web.NewService(
		web.Address(":8081"),
		web.Handler(router),
	)
	service.Run()
}
