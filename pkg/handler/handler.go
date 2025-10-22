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

	books := router.Group("/books", h.userIdentity)
	{

		books.POST("/", h.createBook)
		books.GET("/", h.getAllBook)
		books.GET("/:id", h.getBookById)
		books.PUT("/:id", h.updateBook)
		books.DELETE("/:id", h.deleteBook)

	}
	return router
}
