package middleware

import (
	"gin-photo-storage/constant"
	"gin-photo-storage/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// A wrapper function which returns the auth middleware.
func GetAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		jwtString, err := context.Cookie(constant.JWT)
		if err != nil {
			log.Fatalln(err)
			context.JSON(http.StatusBadRequest, gin.H{
				"code": constant.JWT_MISSING_ERROR,
				"data": make(map[string]string),
				"msg": constant.GetMessage(constant.JWT_MISSING_ERROR),
			})
			context.Abort()
			return
		}

		claim, err := utils.ParseJWT(jwtString)
		if err != nil {
			log.Fatalln(err)
			context.JSON(http.StatusBadRequest, gin.H{
				"code": constant.JWT_PARSE_ERROR,
				"data": make(map[string]string),
				"msg": constant.GetMessage(constant.JWT_PARSE_ERROR),
			})
			context.Abort()
			return
		}

		if utils.IsAuthInRedis(claim.UserName) {
			context.Set("user_name", claim.UserName)
			context.Next()
		} else {
			context.JSON(http.StatusBadRequest, gin.H{
				"code": constant.USER_AUTH_TIMEOUT,
				"data": make(map[string]string),
				"msg": constant.GetMessage(constant.USER_AUTH_TIMEOUT),
			})
			context.Abort()
		}
	}
}
