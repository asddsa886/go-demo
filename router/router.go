package router

import (
	"go-demo/controllers"
	"go-demo/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", controllers.Ping)
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}

	r.Use(middlewares.AuthMiddleWare())
	article := r.Group("/api/article")
	{
		article.POST("/create", controllers.CreateArticle)
		article.GET("/all", controllers.GetAllArticle)
		article.GET("/:id", controllers.GetArticleById)
	}
	return r
}
