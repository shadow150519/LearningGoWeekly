package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hello/WeeklyLearning/jwt-go/demo/model"
	"net/http"
)

func Welcome(ctx *gin.Context)  {
	claims , exist := ctx.Get("claims")
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,"卧槽没拿到绑定的数据")
	}
	token := ctx.GetString("jwtToken")
	if  cusClaims, ok :=  claims.(*model.CustomClaim); ok {
		welcomeStr := fmt.Sprintf("%s,welcome, your token is %s",cusClaims.Username,token)
		ctx.JSON(http.StatusOK,welcomeStr)
	} else {
		ctx.JSON(http.StatusOK,gin.H{
			"msg":"没有收到任何内容111",
		})
	}
}
