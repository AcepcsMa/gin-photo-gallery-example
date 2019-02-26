package models

import (
	"bufio"
	"gin-photo-storage/constant"
	"gin-photo-storage/utils"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"log"
	"mime/multipart"
)

var NoSuchPhotoError = errors.New("no such photo")
var PhotoExistsError = errors.New("photo already exists")
var PhotoFileBrokenError = errors.New("photo file is broken")

// The photo model.
type Photo struct {
	BaseModel
	AuthID 		uint		`json:"auth_id" gorm:"type:int" form:"auth_id"`
	BucketID 	uint		`json:"bucket_id" gorm:"type:int" form:"bucket_id"`
	Name 		string		`json:"name" gorm:"type:varchar(255)" form:"name"`
	Tag 		string		`json:"tag" gorm:"type:varchar(255)" form:"tag"`
	Tags 		[]string	`json:"tags" gorm:"-" form:"tags"`
	Url 		string		`json:"url" gorm:"type:varchar(255)" form:"url"`
	Description string		`json:"description" gorm:"type:text" form:"description"`
	State 		int 		`json:"state" gorm:"type:tinyint(1)" form:"state"`
}

// Add a new photo
func AddPhoto(photoToAdd *Photo, photoFileHeader *multipart.FileHeader) (*Photo, string, error) {
	trx := db.Begin()
	defer trx.Commit()

	// check if the photo exists, select with a WRITE LOCK
	photo := Photo{}
	trx.Set("gorm:query_option", "FOR UPDATE").
		Where("bucket_id = ? AND name = ?", photoToAdd.BucketID, photoToAdd.Name).
		First(&photo)
	if photo.ID > 0 {
		return nil, "", PhotoExistsError
	}

	photo.AuthID = photoToAdd.AuthID
	photo.BucketID = photoToAdd.BucketID
	photo.Name = photoToAdd.Name
	photo.Tag = photoToAdd.Tag
	photo.Description = photoToAdd.Description
	photo.State = 1

	err := trx.Create(&photo).Error
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	err = trx.Model(&Bucket{}).Where("id = ?", photoToAdd.BucketID).
		Update("size", gorm.Expr("size + ?", 1)).
		Error
	if err != nil {
		trx.Rollback()
		log.Println(err)
		return nil, "", err
	}

	// upload to the tencent cloud COS
	if photoFile, err := photoFileHeader.Open(); err == nil {
		uploadID := utils.Upload(photo.ID, photo.Name, bufio.NewReader(photoFile), int(photoFileHeader.Size))
		return &photo, uploadID, nil
	} else {
		return nil, "", PhotoFileBrokenError
	}
}

// Delete a photo by photo id.
func DeletePhotoByID(photoID uint) error {
	trx := db.Begin()
	defer trx.Commit()

	result := trx.Where("id = ? AND state = ?", photoID, 1).Delete(Photo{})
	if err := result.Error; err != nil {
		log.Println(err)
		return err
	}
	if affected := result.RowsAffected; affected == 0 {
		return NoSuchPhotoError
	}
	return nil
}

// Delete a photo by its bucket id & its name.
func DeletePhotoByBucketAndName(bucketID uint, name string) error {
	trx := db.Begin()
	defer trx.Commit()

	result := trx.Where("bucket_id = ? AND name = ?", bucketID, name).Delete(Photo{})
	if err := result.Error; err != nil {
		return err
	}
	if affected := result.RowsAffected; affected == 0 {
		return NoSuchPhotoError
	}
	return nil
}

// Update a photo.
func UpdatePhoto(photoToUpdate *Photo) error {
	trx := db.Begin()
	defer trx.Commit()

	photo := Photo{}
	photo.ID = photoToUpdate.ID

	result := trx.Model(&photo).Updates(photoToUpdate)
	if err := result.Error; err != nil {
		log.Println(err)
		return err
	}
	if affected := result.RowsAffected; affected == 0 {
		return NoSuchPhotoError
	}
	// TODO: update ES

	return nil
}

// Update the url for a photo.
func UpdatePhotoUrl(photoID uint, url string) error {
	trx := db.Begin()
	defer trx.Commit()

	photo := Photo{}
	err := trx.Model(&photo).Update("url", url).Error
	if err != nil {
		return err
	}
	return nil
}

// Get a photo by its photo id.
func GetPhotoByID(photoID uint) (*Photo, error) {
	trx := db.Begin()
	defer trx.Commit()

	photo := Photo{}
	err := trx.Where("id = ?", photoID).First(&photo).Error
	found := NoSuchPhotoError
	if err != nil || photo.ID == 0 {
		log.Println(err)
		found = err
	}
	found = nil
	return &photo, found
}

// Get photos by bucket id.
func GetPhotoByBucketID(bucketID uint, offset int) ([]Photo, error) {
	trx := db.Begin()
	defer trx.Commit()

	photos := make([]Photo, 0, constant.PAGE_SIZE)
	err := trx.Where("bucket_id = ?", bucketID).
		Offset(offset).
		Limit(constant.PAGE_SIZE).
		Find(&photos).
		Error
	if err != nil {
		return photos, err
	}
	return photos, nil
}

// Check photo upload status.
func GetPhotoUploadStatus(uploadID string) int {
	status := utils.GetUploadStatus(uploadID)
	switch status {
	case -2:
		return constant.PHOTO_NOT_EXIST
	case -1:
		return constant.PHOTO_UPLOAD_ERROR
	case 0:
		return constant.PHOTO_UPLOAD_SUCCESS
	case 1:
		return constant.PHOTO_ADD_IN_PROCESS
	default:
		return constant.INVALID_PARAMS
	}
}

func SearchPhotoByTag(tag string, authID int) []Photo {
	return nil
}

func SearchPhotoByDesc(desc string, authID int) []Photo {
	return nil
}