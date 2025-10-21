package handler

import (
	"rest/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)

	}

	api := router.Group("/api") /*, h.userIdentity)*/
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createBook)
			lists.GET("/", h.getAllBook)
			lists.GET("/:id", h.getBookById)
			lists.PUT("/:id", h.updateBook)
			lists.DELETE("/:id", h.deleteBook)
		}
	}
	return router
}
