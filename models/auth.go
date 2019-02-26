package models

import (
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

type Auth struct {
	BaseModel
	UserName 	string `json:"user_name" gorm:"type:varchar(16)"`
	Password 	string `json:"password" gorm:"type:varchar(16)"`
	Email 		string `json:"email" gorm:"type:varchar(128)"`
}

var AuthExistsError = errors.New("auth already exists")

// Add a new auth.
func AddAuth(username, password, email string) error {
	trx := db.Begin()
	defer trx.Commit()

	auth := Auth{}
	trx.Set("gorm:query_option", "FOR UPDATE").
		Where("user_name = ?", username).
		First(&auth)
	if auth.ID > 0 {
		return AuthExistsError
	}

	hash := md5.New()
	io.WriteString(hash, password)	// for safety, don't just save the plain text
	auth.UserName = username
	auth.Password = fmt.Sprintf("%x", hash.Sum(nil))
	auth.Email = email
	err := trx.Create(&auth).Error
	if err != nil {
		return err
	}
	return nil
}

// Check if the auth is valid.
func CheckAuth(username, password string) bool {
	trx := db.Begin()
	defer trx.Commit()

	hash := md5.New()
	io.WriteString(hash, password)
	password = fmt.Sprintf("%x", hash.Sum(nil))	//	for safety, don't just save the plain text
	auth := Auth{}
	trx.Set("gorm:query_option", "FOR UPDATE").
		Where("user_name = ? AND password = ?", username, password).
		First(&auth)
	if auth.ID > 0 {
		return true
	}
	return false
}