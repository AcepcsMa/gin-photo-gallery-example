package models

import (
	"gin-photo-storage/constant"
	"gin-photo-storage/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// The bucket model.
type Bucket struct {
	BaseModel
	AuthID 		uint	`json:"auth_id" gorm:"type:int" form:"auth_id"`
	Name 		string	`json:"bucket_name" gorm:"type:varchar(64)" form:"bucket_name"`
	State 		int		`json:"state" gorm:"type:tinyint(1)" form:"state"`
	Size 		int		`json:"size" gorm:"type:int" form:"bucket_size"`
	Description string	`json:"description" gorm:"type:text" form:"description"`
}

var BucketExistsError = errors.New("bucket already exists")
var NoSuchBucketError = errors.New("no such bucket")

// Add a new bucket.
func AddBucket(bucketToAdd *Bucket) error {
	trx := db.Begin()
	defer trx.Commit()

	// check if the bucket exists, select with a WRITE LOCK.
	bucket := Bucket{}
	trx.Set("gorm:query_option", "FOR UPDATE").
		Where("auth_id = ? AND name = ? AND state = ?", bucketToAdd.AuthID, bucketToAdd.Name, 1).
		First(&bucket)
	if bucket.ID > 0 {
		return BucketExistsError
	}

	bucket.AuthID = bucketToAdd.AuthID
	bucket.Name = bucketToAdd.Name
	bucket.State = 1
	bucket.Size = 0
	bucket.Description = bucketToAdd.Description
	if err := trx.Create(&bucket).Error; err != nil {
		//log.Println(err)
		utils.AppLogger.Info(err.Error(), zap.String("service", "AddBucket()"))
		return err
	}
	return nil
}

// Delete an existed bucket.
func DeleteBucket(bucketID uint) error {
	trx := db.Begin()
	defer trx.Commit()

	result := trx.Where("id = ? and state = ?", bucketID, 1).Delete(Bucket{})
	if err := result.Error; err != nil {
		return err
	}
	if affected := result.RowsAffected; affected == 0 {
		return NoSuchBucketError
	}
	return nil
}

// Update an existed bucket.
func UpdateBucket(bucketToUpdate *Bucket) error {
	trx := db.Begin()
	defer trx.Commit()

	bucket := Bucket{}
	bucket.ID = bucketToUpdate.ID
	result := trx.Model(&bucket).Updates(*bucketToUpdate)
	if err := result.Error; err != nil {
		return err
	}
	if affected := result.RowsAffected; affected == 0 {
		return NoSuchBucketError
	}
	return nil
}

// Get a bucket by bucket id.
func GetBucketByID(bucketID uint) (Bucket, error) {
	trx := db.Begin()
	defer trx.Commit()

	bucket := Bucket{}
	found := NoSuchBucketError
	trx.Where("id = ?", bucketID).First(&bucket)
	if bucket.ID > 0 {
		found = nil
	}
	return bucket, found
}

// Get all buckets of the given user.
func GetBucketByAuthID(authID uint, offset int) ([]Bucket, error) {
	trx := db.Begin()
	defer trx.Commit()

	buckets := make([]Bucket, 0, constant.PAGE_SIZE)
	err := trx.Where("auth_id = ?", authID).
		Offset(offset).
		Limit(constant.PAGE_SIZE).
		Find(&buckets).Error

	if err != nil {
		//log.Println(err)
		utils.AppLogger.Info(err.Error(), zap.String("service", "GetBucketByAuthID()"))
		return buckets, err
	}
	return buckets, nil
}