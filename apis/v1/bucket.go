package v1

import (
	"gin-photo-storage/constant"
	"gin-photo-storage/models"
	"gin-photo-storage/utils"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// Add a new bucket.
func AddBucket(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	bucketToAdd := models.Bucket{}
	if err := context.ShouldBindWith(&bucketToAdd, binding.Form); err != nil {
		//log.Println(err)
		utils.AppLogger.Info(err.Error(), zap.String("service", "AddBucket()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(bucketToAdd.AuthID, "auth_id").Message("Must have auth id")
	validCheck.Required(bucketToAdd.Name, "bucket_name").Message("Must have bucket name")
	validCheck.MaxSize(bucketToAdd.Name, 64, "bucket_name").Message("Bucket name length can not exceed 64")

	if !validCheck.HasErrors() {
		if err := models.AddBucket(&bucketToAdd); err != nil {
			if err == models.BucketExistsError {
				responseCode = constant.BUCKET_ALREADY_EXIST
			} else {
				responseCode = constant.INTERNAL_SERVER_ERROR
			}
		} else {
			responseCode = constant.BUCKET_ADD_SUCCESS
		}
	} else {
		for _, err := range validCheck.Errors {
			//log.Println(err.Message)
			utils.AppLogger.Info(err.Message, zap.String("service", "AddBucket()"))
		}
	}

	data := make(map[string]string)
	data["bucket_name"] = bucketToAdd.Name

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// Delete an existed bucket.
func DeleteBucket(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	bucketID, bucketErr := strconv.Atoi(context.Query("bucket_id"))
	if bucketErr != nil {
		//log.Println(bucketErr)
		utils.AppLogger.Info(bucketErr.Error(), zap.String("service", "DeleteBucket()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(bucketID, "bucket_id").Message("Must have bucket id")
	validCheck.Min(bucketID, 1, "bucket_id").Message("Bucket id should be positive")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if err := models.DeleteBucket(uint(bucketID)); err != nil {
			if err == models.NoSuchBucketError {
				responseCode = constant.BUCKET_NOT_EXIST
			} else {
				responseCode = constant.INTERNAL_SERVER_ERROR
			}
		} else {
			responseCode = constant.BUCKET_DELETE_SUCCESS
		}
	} else {
		for _, err := range validCheck.Errors {
			//log.Println(err.Message)
			utils.AppLogger.Info(err.Message, zap.String("service", "DeleteBucket()"))
		}
	}

	data["bucket_id"] = bucketID
	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// Update an existed bucket.
func UpdateBucket(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS

	bucketToUpdate := models.Bucket{}
	if err := context.ShouldBindWith(&bucketToUpdate, binding.Form); err != nil {
		//log.Println(err)
		utils.AppLogger.Info(err.Error(), zap.String("service", "UpdateBucket()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(bucketToUpdate.ID, "bucket_id").Message("Must have bucket id")
	validCheck.MaxSize(bucketToUpdate.Name, 64, "bucket_name").Message("Bucket name length can not exceed 64")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if err := models.UpdateBucket(&bucketToUpdate); err != nil {
			if err == models.NoSuchBucketError {
				responseCode = constant.BUCKET_NOT_EXIST
			} else {
				responseCode = constant.INTERNAL_SERVER_ERROR
			}
		} else {
			responseCode = constant.BUCKET_UPDATE_SUCCESS
		}
	} else {
		for _, err := range validCheck.Errors {
			//log.Println(err.Message)
			utils.AppLogger.Info(err.Message, zap.String("service", "UpdateBucket()"))
		}
	}

	data["bucket_id"] = bucketToUpdate.ID
	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// Get a bucket by bucket id.
func GetBucketByID(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	bucketID, bucketErr := strconv.Atoi(context.Query("bucket_id"))
	if bucketErr != nil {
		//log.Println(bucketErr)
		utils.AppLogger.Info(bucketErr.Error(), zap.String("service", "GetBucketByID()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(bucketID, "bucket_id").Message("Must have bucket id")
	validCheck.Min(bucketID, 1, "bucket_id").Message("Bucket id should be positive")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if bucket, err := models.GetBucketByID(uint(bucketID)); err != nil {
			if err == models.NoSuchBucketError {
				responseCode = constant.BUCKET_NOT_EXIST
			} else {
				responseCode = constant.INTERNAL_SERVER_ERROR
			}
		} else {
			responseCode = constant.BUCKET_GET_SUCCESS
			data["bucket"] = bucket
		}
	} else {
		for _, err := range validCheck.Errors {
			//log.Println(err.Message)
			utils.AppLogger.Info(err.Message, zap.String("service", "GetBucketByID()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// Get buckets by auth id.
func GetBucketByAuthID(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	authID, authErr := strconv.Atoi(context.Query("auth_id"))
	offset := context.GetInt("offset")
	if authErr != nil{
		//log.Println(authErr)
		utils.AppLogger.Info(authErr.Error(), zap.String("service", "GetBucketByAuthID()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(authID, "auth_id").Message("Must have auth id")
	validCheck.Min(authID, 1, "auth_id").Message("Auth id should be positive")
	validCheck.Min(offset, 0, "page_offset").Message("Page offset must be >= 0")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if buckets, err := models.GetBucketByAuthID(uint(authID), offset); err != nil {
			responseCode = constant.INTERNAL_SERVER_ERROR
		} else {
			responseCode = constant.BUCKET_GET_SUCCESS
			data["buckets"] = buckets
		}
	} else {
		for _, err := range validCheck.Errors {
			//log.Println(err.Message)
			utils.AppLogger.Info(err.Message, zap.String("service", "GetBucketByAuthID()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}