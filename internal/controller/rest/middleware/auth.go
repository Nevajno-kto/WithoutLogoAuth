package middleware

import (
	"net/http"
	"strings"

	"github.com/nevajno-kto/without-logo-auth/config"
	authjwt "github.com/nevajno-kto/without-logo-auth/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}

		if err := c.ShouldBindHeader(&h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				c.JSON(http.StatusBadRequest, gin.H{
					"ok":        false,
					"error_msg": invalidArgs,
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"ok":        false,
				"error_msg": nil,
			})
			c.Abort()
			return
		}

		idTokenHeader := strings.Split(h.IDToken, "Bearer ")

		if len(idTokenHeader) < 2 {

			c.JSON(http.StatusBadRequest, gin.H{
				"ok":        false,
				"error_msg": "Must provide Authorization header with format `Bearer {token}`",
			})
			c.Abort()
			return
		}

		userid, err := authjwt.ParseToken(idTokenHeader[1], []byte(config.GetConfig().JWT.Secret))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":        false,
				"error_msg": "Provided token is invalid",
			})
			c.Abort()
			return
		}

		c.Set("userid", userid)

		c.Next()
	}
}
