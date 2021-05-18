package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"hello/WeeklyLearning/jwt-go/demo/model"
	"net/http"
)

var SecretKey = []byte("wuhuqifei")

func ValidateWelcome(ctx *gin.Context)  {

	tokenStr := ctx.GetHeader("auth")
	if tokenStr == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,gin.H{
			"msg":"未提供token，请登录",
		})
	}
		token, err := jwt.ParseWithClaims(tokenStr,&model.CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest,gin.H{
				"msg": err.Error(),
			})
		}
		if customClaim ,ok := token.Claims.(*model.CustomClaim); ok {
			ctx.Set("jwtToken",tokenStr)
			ctx.Set("claims",customClaim)
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest,gin.H{
				"msg": customClaim.Valid().Error(),
			})
		}


}