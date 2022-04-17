package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/nevajno-kto/without-logo-auth/internal/usecase"
	"github.com/nevajno-kto/without-logo-auth/pkg/logger"
)

func NewRouter(handler *gin.Engine, l logger.Interface, a *usecase.AuthUseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.StaticFile("/code", "./code.txt")

	h := handler.Group("/clients")
	{
		newAuthRoutes(h, l, *a)
	}
}
