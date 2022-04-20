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
	Action       string `json:"action" binding:"required"`
	Organization string `json:"org" binding:"required"`
	Phone        string `json:"phone" binding:"required,e164"`
	UserType     string `json:"userType" binding:"omitempty,min=5,max=6"`
	Password     string `json:"password" binding:"omitempty,gte=8,lte=50"`
	Code         int    `json:"code" binding:"omitempty,gte=1000,lte=9999"`
}

func (r *ClientsRoutes) SignUp(c *gin.Context) {
	var action, userType int
	var request authRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "rest - SignUp - invalid request body")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	switch request.Action {
	case "request":
		action = entity.SignUpRequest
	case "confirm":
		action = entity.SignUpConfirm
	default:
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	switch request.UserType {
	case "client":
		userType = entity.Client
	case "admin":
		userType = entity.Admin
	default:
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	authResult, err := r.u.SignUp(
		c.Request.Context(),
		entity.Auth{
			Action: action,
			User: entity.User{
				Phone:        request.Phone,
				Password:     request.Password,
				Organization: request.Organization,
				Type:         userType,
			},
			Code: request.Code,
		})

	if err != nil {
		if errors.Is(err, entity.ErrServiceProblem) {
			r.l.Error(errors.Unwrap(err), "rest - SignUp - signup service problems")
			errorResponse(c, http.StatusInternalServerError, "internal server error")
			return
		}

		if errors.Is(err, entity.ErrTimeout) {
			timeoutResponse(c, http.StatusOK, authResult)
			return
		}

		errorResponse(c, http.StatusOK, err.Error())
		return
	}

	response(c, http.StatusOK, authResult)
}

func (r *ClientsRoutes) SignIn(c *gin.Context) {
	var action int
	var request authRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "rest - SignIn - invalid request body")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	switch request.Action {
	case "password":
		action = entity.SignInByPassword
	case "request":
		action = entity.SignInRequest
	case "confirm":
		action = entity.SignInConfirm
	default:
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// switch request.UserType {
	// case "client":
	// 	userType = entity.Client
	// case "admin":
	// 	userType = entity.Admin
	// default:
	// 	errorResponse(c, http.StatusBadRequest, "invalid request body")
	// 	return
	// }

	authResult, err := r.u.SignIn(
		c.Request.Context(),
		entity.Auth{
			Action: action,
			User: entity.User{
				Phone:        request.Phone,
				Password:     request.Password,
				Organization: request.Organization,
			},
			Code: request.Code,
		})

	if err != nil {

		if errors.Is(err, entity.ErrServiceProblem) {
			r.l.Error(err, "rest - SignIn - signup service problems")
			errorResponse(c, http.StatusInternalServerError, "internal server error")
			return
		}

		if errors.Is(err, entity.ErrTimeout) {
			timeoutResponse(c, http.StatusOK, authResult)
			return
		}

		errorResponse(c, http.StatusOK, err.Error())
		return
	}

	response(c, http.StatusOK, authResult)

}

func (r *ClientsRoutes) Test(c *gin.Context) {
	response(c, http.StatusOK, "SUCCESS")
}
