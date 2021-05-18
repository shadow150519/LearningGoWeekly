package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"hello/WeeklyLearning/jwt-go/demo/middleware"
	"hello/WeeklyLearning/jwt-go/demo/model"
	"net/http"
)

func Login(ctx *gin.Context)  {
	UserInfo := struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Age int `json:"age" binding:"required"`
	}{}
	UserInfo.Username = ctx.PostForm("username")
	UserInfo.Password = ctx.PostForm("password")
	// 这里就直接写死了...
	if UserInfo.Username == "admin" && UserInfo.Password == "123456" {
		// 登录成功，返回token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,model.CustomClaim{
			Username:       "admin",
			Age:            18,
			StandardClaims: jwt.StandardClaims{},
		})
		tokenStr, err := token.SignedString(middleware.SecretKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,"生成token时出错")
			return
		}
		ctx.Header("auth",fmt.Sprintf("bear %s", tokenStr))
		ctx.JSON(http.StatusOK,"登陆成功")
	} else {
		ctx.JSON(http.StatusOK,"用户名或密码错误")
		return
	}

}