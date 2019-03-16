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
	"strings"
)

// Add a new photo
func AddPhoto(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS

	photoFile, fileErr := context.FormFile("photo")
	if fileErr != nil {
		//log.Println(fileErr)
		utils.AppLogger.Info(fileErr.Error(), zap.String("service", "AddPhoto()"))
	}

	photo := models.Photo{}
	paramErr := context.ShouldBindWith(&photo, binding.Form)

	if fileErr != nil || paramErr != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(photo.AuthID, "auth_id").Message("Must have auth id")
	validCheck.Required(photo.BucketID, "bucket_id").Message("Must have bucket id")
	validCheck.Required(photo.Name, "photo_name").Message("Must have photo name")
	validCheck.MaxSize(photo.Name, 255, "photo_name").Message("Photo name len must not exceed 255")

	data := make(map[string]interface{})
	photoToAdd := &models.Photo{BucketID: photo.BucketID, AuthID: photo.AuthID,
		Name: photo.Name, Description: photo.Description,
		Tag:strings.Join(photo.Tags, ";")}

	if !validCheck.HasErrors() {
		if photoToAdd, uploadID, err := models.AddPhoto(photoToAdd, photoFile); err != nil {
			if err == models.PhotoExistsError {
				responseCode = constant.PHOTO_ALREADY_EXIST
			} else {
				responseCode = constant.INTERNAL_SERVER_ERROR
			}
		} else {
			responseCode = constant.PHOTO_ADD_IN_PROCESS
			data["photo"] = *photoToAdd
			data["photo_upload_id"] = uploadID
		}
	} else {
		for _, err := range validCheck.Errors {
			//log.Println(err.Message)
			utils.AppLogger.Info(err.Message, zap.String("service", "AddPhoto()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// Delete an existed photo.
func DeletePhoto(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	bucketID, err := strconv.Atoi(context.PostForm("bucket_id"))
	photoName := context.PostForm("photo_name")
	if err != nil {
		//log.Println(err)
		utils.AppLogger.Info(err.Error(), zap.String("service", "DeletePhoto()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(bucketID, "bucket_id").Message("Must have bucket id")
	validCheck.Required(photoName, "photo_name").Message("Must have photo name")
	validCheck.MaxSize(photoName, 255, "photo_name").Message("Photo name length must not exceed 255")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if err := models.DeletePhotoByBucketAndName(uint(bucketID), photoName); err != nil {
			if err == models.NoSuchPhotoError {
				responseCode = constant.PHOTO_NOT_EXIST
			} else {
				responseCode = constant.INTERNAL_SERVER_ERROR
			}
		} else {
			responseCode = constant.PHOTO_DELETE_SUCCESS
		}
	} else {
		for _, err := range validCheck.Errors {
			//log.Println(err.Message)
			utils.AppLogger.Info(err.Message, zap.String("service", "DeletePhoto()"))
		}
	}

	data["photo_name"] = photoName
	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// Update an existed photo.
func UpdatePhoto(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	photoToUpdate := models.Photo{}
	err := context.ShouldBindWith(&photoToUpdate, binding.Form)
	if err != nil {
		utils.AppLogger.Info(err.Error(), zap.String("service", "UpdatePhoto()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg": constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(photoToUpdate.ID, "photo_id").Message("Must have photo id")
	validCheck.Min(int(photoToUpdate.ID), 1, "photo_id").Message("Photo id should be positive")
	validCheck.MaxSize(photoToUpdate.Name, 255, "photo_name").
		Message("Photo name length can not exceed 255")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		photoToUpdate.Tag = strings.Join(photoToUpdate.Tags, ";")
		if photo, err := models.UpdatePhoto(&photoToUpdate); err != nil {
			if err == models.NoSuchPhotoError {
				responseCode = constant.PHOTO_NOT_EXIST
			} else {
				responseCode = constant.INTERNAL_SERVER_ERROR
			}
		} else {
			responseCode = constant.PHOTO_UPDATE_SUCCESS
			data["photo"] = *photo
		}
	} else {
		for _, err := range validCheck.Errors {
			//log.Println(err.Message)
			utils.AppLogger.Info(err.Message, zap.String("service", "UpdatePhoto()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// Get a photo by photo id.
func GetPhotoByID(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	photoID, photoErr := strconv.Atoi(context.Query("photo_id"))
	if photoErr != nil {
		//log.Println(photoErr)
		utils.AppLogger.Info(photoErr.Error(), zap.String("service", "GetPhotoByID()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(photoID, "photo_id").Message("Must have photo id")
	validCheck.Min(photoID, 1, "photo_id").Message("Photo id should be positive")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if photo, err := models.GetPhotoByID(uint(photoID)); err != nil {
			if err == models.NoSuchPhotoError {
				responseCode = constant.PHOTO_NOT_EXIST
			} else {
				responseCode = constant.INTERNAL_SERVER_ERROR
			}
		} else {
			responseCode = constant.PHOTO_GET_SUCCESS
			photo.Tags = strings.Split(photo.Tag, ";")
			data["photo"] = *photo
		}
	} else {
		for _, err := range validCheck.Errors {
			//log.Println(err.Message)
			utils.AppLogger.Info(err.Message, zap.String("service", "GetPhotoByID()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// Get photos by bucket id.
func GetPhotoByBucketID(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	bucketID, bucketErr := strconv.Atoi(context.Query("bucket_id"))
	offset := context.GetInt("offset")
	if bucketErr != nil{
		utils.AppLogger.Info(bucketErr.Error(), zap.String("service", "GetPhotoByBucketID()"))
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
	validCheck.Min(offset, 0, "page_offset").Message("Page offset must be >= 0")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if photos, err := models.GetPhotoByBucketID(uint(bucketID), offset); err != nil {
			responseCode = constant.INTERNAL_SERVER_ERROR
		} else {
			responseCode = constant.PHOTO_GET_SUCCESS
			for i := 0;i < len(photos);i++ {
				photos[i].Tags = strings.Split(photos[i].Tag, ";")
			}
			data["photo"] = photos
		}
	} else {
		for _, err := range validCheck.Errors {
			//log.Println(err.Message)
			utils.AppLogger.Info(err.Message, zap.String("service", "GetPhotoByBucketID()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// Get the upload status of a photo by upload id.
func GetPhotoUploadStatus(context *gin.Context) {
	uploadID := context.Query("upload_id")

	validCheck := validation.Validation{}
	validCheck.Required(uploadID, "upload_id").Message("Must have upload id")

	responseCode := constant.INVALID_PARAMS
	data := make(map[string]interface{})
	data["upload_id"] = uploadID
	if !validCheck.HasErrors() {
		responseCode = models.GetPhotoUploadStatus(uploadID)
	} else {
		for _, err := range validCheck.Errors {
			//log.Println(err)
			utils.AppLogger.Info(err.Message, zap.String("service", "GetPhotoUploadStatus()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg": constant.GetMessage(responseCode),
	})
}

// Search a photo (by tag / description)
func SearchPhoto(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	authID, err := strconv.Atoi(context.Query("auth_id"))
	tag, tagExisted := context.GetQuery("tag")
	desc, descExisted := context.GetQuery("desc")
	if err != nil || (tagExisted && descExisted) || (!tagExisted && !descExisted) {
		utils.AppLogger.Info(constant.GetMessage(responseCode), zap.String("service", "SearchPhoto()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
	}

	var searchType models.SearchType
	var field string
	if tagExisted {
		searchType = constant.SEARCH_BY_TAG
		field = tag
	} else {
		searchType = constant.SEARCH_BY_DESC
		field = desc
	}

	validCheck := validation.Validation{}
	validCheck.Min(authID, 1, "auth_id").Message("Auth id must be positive")
	validCheck.MinSize(field, 1, "search_field").Message("Search field can't be empty")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		offset := context.GetInt("offset")
		if photos, err := models.SearchPhoto(field, uint(authID), offset, searchType); err == nil {
			data["photos"] = photos
			responseCode = constant.PHOTO_SEARCH_BY_TAG_SUCCESS
		} else {
			responseCode = constant.INTERNAL_SERVER_ERROR
		}
	} else {
		for _, err := range validCheck.Errors {
			utils.AppLogger.Info(err.Message, zap.String("service", "SearchPhoto()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg": constant.GetMessage(responseCode),
	})
}