package v1

import (
	"gin-photo-storage/conf"
	"gin-photo-storage/models"
	"gin-photo-storage/constant"
	"gin-photo-storage/utils"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Add a new auth.
func AddAuth(context *gin.Context) {

	userName := context.PostForm("user_name")
	password := context.PostForm("password")
	email := context.PostForm("email")

	// set up param validation
	validCheck := validation.Validation{}
	validCheck.Required(userName, "user_name").Message("Must have user name")
	validCheck.MaxSize(userName, 16, "user_name").Message("User name length can not exceed 16")
	validCheck.MinSize(userName, 6, "user_name").Message("User name length is at least 6")
	validCheck.Required(password, "password").Message("Must have password")
	validCheck.MaxSize(password, 16, "password").Message("Password length can not exceed 16")
	validCheck.MinSize(password, 6, "password").Message("Password length is at least 6")
	validCheck.Required(email, "email").Message("Must have email")
	validCheck.MaxSize(email, 128, "email").Message("Email can not exceed 128 chars")

	responseCode := constant.INVALID_PARAMS
	if !validCheck.HasErrors() {
		if err := models.AddAuth(userName, password, email); err == nil {
			responseCode = constant.USER_ADD_SUCCESS
		} else {
			responseCode = constant.USER_ALREADY_EXIST
		}
	} else {
		for _, err := range validCheck.Errors {
			log.Println(err)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": userName,
		"msg": constant.GetMessage(responseCode),
	})
}

// Check if an auth is valid.
func CheckAuth(context *gin.Context) {

	userName := context.PostForm("user_name")
	password := context.PostForm("password")

	// set up param validation
	validCheck := validation.Validation{}
	validCheck.Required(userName, "user_name").Message("Must have user name")
	validCheck.MaxSize(userName, 16, "user_name").Message("User name length can not exceed 16")
	validCheck.MinSize(userName, 6, "user_name").Message("User name length is at least 6")
	validCheck.Required(password, "password").Message("Must have password")
	validCheck.MaxSize(password, 16, "password").Message("Password length can not exceed 16")
	validCheck.MinSize(password, 6, "password").Message("Password length is at least 6")

	responseCode := constant.INVALID_PARAMS
	if !validCheck.HasErrors() {
		if models.CheckAuth(userName, password) {
			if jwtString, err := utils.GenerateJWT(userName); err != nil {
				responseCode = constant.JWT_GENERATION_ERROR
			} else {
				// pass auth validation
				// 1. set JWT to user's cookie
				// 2. add user to the Redis
				context.SetCookie(constant.JWT, jwtString,
					constant.COOKIE_MAX_AGE, conf.ServerCfg.Get(constant.SERVER_PATH),
					conf.ServerCfg.Get(constant.SERVER_DOMAIN), true, true)
				if err = utils.AddAuthToRedis(userName); err != nil {
					responseCode = constant.INTERNAL_SERVER_ERROR
				} else {
					responseCode = constant.USER_AUTH_SUCCESS
				}
			}
		} else {
			responseCode = constant.USER_AUTH_ERROR
		}
	} else {
		for _, err := range validCheck.Errors {
			log.Println(err)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": userName,
		"msg": constant.GetMessage(responseCode),
	})
}

//func SignOut(context *gin.Context) {
//	userName, _ := context.Get("user_name")
//	data := make(map[string]string)
//	data["user_name"] = userName.(string)
//	responseCode := constant.INVALID_PARAMS
//	if utils.RemoveAuthFromRedis(userName.(string)) {
//		responseCode = constant.USER_SIGNOUT_SUCCESS
//	} else {
//		responseCode = constant.INTERNAL_SERVER_ERROR
//	}
//
//	context.JSON(http.StatusOK, gin.H{
//		"code": responseCode,
//		"data": data,
//		"msg": constant.GetMessage(responseCode),
//	})
//}