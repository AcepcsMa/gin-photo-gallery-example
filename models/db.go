package models

import (
	"errors"
	"fmt"
	"gin-photo-storage/conf"
	"gin-photo-storage/constant"
	"gin-photo-storage/utils"
	_ "github.com/go-sql-driver/mysql" // remember to import mysql driver
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

var db *gorm.DB
var CallbackUpdateError = errors.New("callback update error")

// The base model of all models, including ID & CreatedAt & UpdatedAt.
type BaseModel struct {
	ID 			uint 		`json:"id" gorm:"primary_key;AUTO_INCREMENT" form:"id"`
	CreatedAt 	time.Time 	`json:"created_at" gorm:"default: CURRENT_TIMESTAMP" form:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at" gorm:"default: CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" form:"updated_at"`
}

// Init the database connection.
func init() {
	dbType := conf.ServerCfg.Get(constant.DB_TYPE)
	dbHost := conf.ServerCfg.Get(constant.DB_HOST)
	dbPort := conf.ServerCfg.Get(constant.DB_PORT)
	dbUser := conf.ServerCfg.Get(constant.DB_USER)
	dbPwd := conf.ServerCfg.Get(constant.DB_PWD)
	dbName := conf.ServerCfg.Get(constant.DB_NAME)

	var err error
	db, err = gorm.Open(dbType, fmt.Sprintf(constant.DB_CONNECT, dbUser, dbPwd, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatalln("Fail to connect database!")
	}

	db.SingularTable(true)
	if !db.HasTable(&Auth{}) {
		db.CreateTable(&Auth{})
	}
	if !db.HasTable(&Bucket{}) {
		db.CreateTable(&Bucket{})
	}
	if !db.HasTable(&Photo{}) {
		db.CreateTable(&Photo{})
	}

	go ListenRedisCallback()	// launch a background goroutine to listen to callbacks from redis
}

// Listen to callback messages from redis channels.
// 1. When a photo is uploaded successfully, the callback asks to update the photo url in the db.
// 2. When it fails to upload a photo, the callback asks to delete the photo record in the db.
func ListenRedisCallback() {

	// wait until utils package is initialized
	<- utils.InitComplete

	// subscribe redis channels
	updateChan := utils.RedisClient.Subscribe(constant.URL_UPDATE_CHANNEL).Channel()
	deleteChan := utils.RedisClient.Subscribe(constant.PHOTO_DELETE_CHANNEL).Channel()

	// loop and listen
	for {
		select {
		case msg := <-updateChan:
			photoID, _ := strconv.Atoi(msg.Payload[:strings.Index(msg.Payload, "-")])
			photoUrl := msg.Payload[strings.Index(msg.Payload, "-") + 1:]
			dbErr := UpdatePhotoUrl(uint(photoID), photoUrl)
			esErr := AddPhotoUrl(uint(photoID), photoUrl)
			if dbErr != nil || esErr != nil {
				log.Println(CallbackUpdateError)
			} else {
				utils.SetUploadStatus(fmt.Sprintf(constant.PHOTO_UPDATE_ID_FORMAT, photoID), 0)
			}
			//if err := UpdatePhotoUrl(uint(photoID), photoUrl); err != nil {
			//	log.Println(err)
			//} else {
			//	if err := AddPhotoUrl(uint(photoID), photoUrl); err != nil {
			//		log.Println(err)
			//	} else {
			//		utils.SetUploadStatus(fmt.Sprintf(constant.PHOTO_UPDATE_ID_FORMAT, photoID), 0)
			//	}
			//}
		case msg := <- deleteChan:
			photoID, _ := strconv.Atoi(msg.Payload)
			if err := DeletePhotoByID(uint(photoID)); err != nil {
				log.Println(err)
			} else {
				utils.SetUploadStatus(fmt.Sprintf(constant.PHOTO_UPDATE_ID_FORMAT, photoID), -1)
			}
		default:
		}
	}
}