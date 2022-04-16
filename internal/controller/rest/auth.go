package rest

import (
	"errors"
	"net/http"

	"github.com/nevajno-kto/without-logo-auth/internal/controller/rest/middleware"
	"github.com/nevajno-kto/without-logo-auth/internal/entity"
	"github.com/nevajno-kto/without-logo-auth/internal/usecase"
	"github.com/nevajno-kto/without-logo-auth/pkg/logger"

	"github.com/gin-gonic/gin"
)

type ClientsRoutes struct {
	u usecase.AuthUseCase
	l logger.Interface
}

func newAuthRoutes(h *gin.RouterGroup, l logger.Interface, a usecase.AuthUseCase) {

	r := ClientsRoutes{u: a, l: l}

	h.POST("/signup", r.SignUp)
	h.POST("/signin", r.SignIn)
	h.Use(middleware.AuthUser())
	h.POST("/testdrive", r.Test)
}

type authRequest struct {
	Type         string `json:"type" binding:"required"`
	Organization string `json:"org" binding:"required"`
	Phone        string `json:"phone" binding:"required,e164"`
	Password     string `json:"password" binding:"omitempty,gte=8,lte=50"`
	Code         int    `json:"code" binding:"omitempty,gte=100000,lte=999999"`
	Name         string `json:"name" binding:"omitempty"`
}

func (r *ClientsRoutes) SignUp(c *gin.Context) {
	var request authRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "rest - SignUp - invalid request body")
		errorResponse(c, http.StatusBadRequest, "Ошибка в теле запроса")
		return
	}

	err := r.u.SignUp(
		c.Request.Context(),
		entity.Auth{
			Type: request.Type,
			User: entity.User{
				Phone:        request.Phone,
				Name:         request.Name,
				Password:     request.Password,
				Organization: request.Organization,
			},
			Code: request.Code,
		})

	if err != nil {
		if errors.Is(err, entity.ErrSignUp) {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		r.l.Error(err, "rest - SignUp - signup service problems")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response(c, http.StatusOK, nil)
}

func (r *ClientsRoutes) SignIn(c *gin.Context) {
	var request authRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "rest - SignIn - invalid request body")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := r.u.SignIn(
		c.Request.Context(),
		entity.Auth{
			Type: request.Type,
			User: entity.User{
				Phone:        request.Phone,
				Name:         request.Name,
				Password:     request.Password,
				Organization: request.Organization,
			},
			Code: request.Code,
		})

	if err != nil {
		if errors.Is(err, entity.ErrSingIn) {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		r.l.Error(err, "rest - SignIn - signup service problems")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response(c, http.StatusOK, token)

}

func (r *ClientsRoutes) Test(c *gin.Context) {
	response(c, http.StatusOK, "SUCCESS")
}
