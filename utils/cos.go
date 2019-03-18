package utils

import (
	"context"
	"fmt"
	"gin-photo-storage/conf"
	"gin-photo-storage/constant"
	"github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
)

var (
	CosClient    *cos.Client
	CosUrlFormat = "http://%s-%s.cos.%s.myqcloud.com"
	BucketName   = ""
	AppID        = ""
	Region       = ""
)

// init COS client
func init() {
	BucketName = conf.ServerCfg.Get(constant.COS_BUCKET_NAME)
	AppID = conf.ServerCfg.Get(constant.COS_APP_ID)
	Region = conf.ServerCfg.Get(constant.COS_REGION)
	u, _ := url.Parse(fmt.Sprintf(CosUrlFormat, BucketName, AppID, Region))
	b := &cos.BaseURL{BucketURL: u}
	CosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.ServerCfg.Get(constant.COS_SECRET_ID),
			SecretKey: conf.ServerCfg.Get(constant.COS_SECRET_KEY),
		},
	})
	//log.Printf("COS client %s init", CosClient.BaseURL)
	AppLogger.Info("COS client init", zap.String("service", "init()"))
}

// upload a photo to the tencent cloud COS
func Upload(photoID uint, fileName string, file io.Reader, fileSize int) string {
	uploadID := fmt.Sprintf(constant.PHOTO_UPDATE_ID_FORMAT, photoID)
	go AsyncUpload(uploadID, photoID, fileName, file, fileSize)	// upload in the ASYNC way
	return uploadID
}

// upload a photo to the tencent cloud COS in the ASYNC way
func AsyncUpload(uploadID string, photoID uint, fileName string, file io.Reader, fileSize int) {
	// set upload status in redis
	if !SetUploadStatus(uploadID, 1) {
		//log.Println("Fail to set upload status before upload.")
		AppLogger.Info("Fail to set upload status before upload.", zap.String("service", "AsyncUpload()"))
		return
	}

	// upload the photo using COS SDK
	putOption := cos.ObjectPutOptions{}
	putOption.ObjectPutHeaderOptions = &cos.ObjectPutHeaderOptions{ContentLength: fileSize}
	_, err := CosClient.Object.Put(context.Background(), fileName, file, &putOption)

	// upload fails, send callback asking for photo deletion
	if err != nil {
		//log.Println(err)
		AppLogger.Info(err.Error(), zap.String("service", "AsyncUpload()"))
		if !SendToChannel(constant.PHOTO_DELETE_CHANNEL, fmt.Sprintf("%d", photoID)) {
			//log.Println("Fail to send delete-photo message to channel")
			AppLogger.Info("Fail to send delete-photo msg to channel.", zap.String("service", "AsyncUpload()"))
		}
		return
	}

	// upload success, send callback asking for updating the photo url
	fileUrl := fmt.Sprintf(CosUrlFormat, BucketName, AppID, Region) + "/" + fileName
	updateUrlMessage := fmt.Sprintf("%d-%s", photoID, fileUrl)
	if !SendToChannel(constant.URL_UPDATE_CHANNEL, updateUrlMessage) {
		//log.Println("Fail to send update-photo-url message to channel")
		AppLogger.Info("Fail to send update-photo-url msg to channel.", zap.String("service", "AsyncUpload()"))
	}
}