package v1

import (
	"gin-photo-storage/models"
	"gin-photo-storage/constant"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func AddPhoto(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS

	paramErr := false
	photoFile, err := context.FormFile("photo")
	if err != nil {
		log.Println(err)
		paramErr = true
	}

	photoName := context.PostForm("photo_name")
	photoTags := context.PostFormArray("photo_tags")
	photoDesc := context.PostForm("description")
	authID, authErr := strconv.Atoi(context.PostForm("auth_id"))
	bucketID, bucketErr := strconv.Atoi(context.PostForm("bucket_id"))
	if authErr != nil || bucketErr != nil {
		if authErr != nil {
			log.Println(authErr)
		}
		if bucketErr != nil {
			log.Println(bucketErr)
		}
		paramErr = true
	}

	if paramErr {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(authID, "auth_id").Message("Must have auth id")
	validCheck.Required(bucketID, "bucket_id").Message("Must have bucket id")
	validCheck.Required(photoName, "photo_name").Message("Must have photo name")
	validCheck.MaxSize(photoName, 255, "photo_name").Message("Photo name len must not exceed 255")

	data := make(map[string]interface{})
	photoToAdd := &models.Photo{BucketID: uint(bucketID), AuthID: uint(authID),
								Name: photoName, Description: photoDesc,
								Tag:strings.Join(photoTags, ";")}
	uploadID := ""
	if !validCheck.HasErrors() {
		if photoToAdd, uploadID, err = models.AddPhoto(photoToAdd, photoFile); err != nil {
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
			log.Println(err.Message)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

func DeletePhoto(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	bucketID, err := strconv.Atoi(context.PostForm("bucket_id"))
	photoName := context.PostForm("photo_name")
	if err != nil {
		log.Println(err)
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
			log.Println(err.Message)
		}
	}

	data["photo_name"] = photoName
	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

func UpdatePhoto(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	photoToUpdate := models.Photo{}
	err := context.ShouldBindWith(&photoToUpdate, binding.Form)
	if err != nil {
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
		if err := models.UpdatePhoto(&photoToUpdate); err != nil {
			if err == models.NoSuchPhotoError {
				responseCode = constant.PHOTO_NOT_EXIST
			} else {
				responseCode = constant.INTERNAL_SERVER_ERROR
			}
		} else {
			responseCode = constant.PHOTO_UPDATE_SUCCESS
		}
	} else {
		for _, err := range validCheck.Errors {
			log.Println(err.Message)
		}
	}

	data["photo"] = photoToUpdate
	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

func GetPhotoByID(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	photoID, photoErr := strconv.Atoi(context.Query("photo_id"))
	if photoErr != nil {
		log.Println(photoErr)
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
			log.Println(err.Message)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

func GetPhotoByBucketID(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	bucketID, bucketErr := strconv.Atoi(context.Query("bucket_id"))
	offset := context.GetInt("offset")
	if bucketErr != nil{
		if bucketErr != nil {
			log.Println(bucketErr)
		}
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
			log.Println(err.Message)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

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
		for err := range validCheck.Errors {
			log.Println(err)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg": constant.GetMessage(responseCode),
	})
}