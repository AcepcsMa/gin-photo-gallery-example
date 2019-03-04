package middleware

import (
	"gin-photo-storage/conf"
	"gin-photo-storage/constant"
	"gin-photo-storage/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// A wrapper function which returns the refresh middleware.
func GetRefreshMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if userName, ok := context.Get("user_name"); ok {
			// generate a new valid JWT for the user
			jwtString, err := utils.GenerateJWT(userName.(string))
			if err != nil {
				log.Fatalln(err)
				data := make(map[string]string)
				data["user_name"] = userName.(string)
				context.JSON(http.StatusBadRequest, gin.H{
					"code": constant.JWT_GENERATION_ERROR,
					"data": data,
					"msg": constant.GetMessage(constant.JWT_GENERATION_ERROR),
				})
				context.Abort()
				return
			}

			// save the new JWT in user's cookie
			context.SetCookie(constant.JWT, jwtString,
				constant.COOKIE_MAX_AGE, "/",
				conf.ServerCfg.Get(constant.SERVER_DOMAIN), true, true)

			// refresh user in the redis
			err = utils.AddAuthToRedis(userName.(string))
			if err != nil {
				log.Fatalln(err)
				context.JSON(http.StatusBadRequest, gin.H{
					"code": constant.INTERNAL_SERVER_ERROR,
					"data": make(map[string]string),
					"msg": constant.GetMessage(constant.INTERNAL_SERVER_ERROR),
				})
				context.Abort()
				return
			}
			context.Next()
		}
		context.Abort()
	}
}
