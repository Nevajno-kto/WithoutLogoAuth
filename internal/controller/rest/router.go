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

	//******************************** DEBUG ************************************
	handler.StaticFile("/signUpcode", "./signUpcode.txt")
	handler.StaticFile("/signIncode", "./signIncode.txt")
	//******************************** DEBUG ************************************

	h := handler.Group("/users")
	{
		newAuthRoutes(h, l, *a)
	}
}
