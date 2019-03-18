package models

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin-photo-storage/conf"
	"gin-photo-storage/constant"
	"gin-photo-storage/utils"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"go.uber.org/zap"
	"strings"
)

var ESClient *elasticsearch.Client
var PhotoIndexingError = errors.New("photo indexing error")
var PhotoSearchError = errors.New("photo search error")
var PhotoUpdateError = errors.New("photo update error")

var SearchRequest = `{
	"query": {
		"bool": {
			"must": [
				{
					"match": {
						"%s": "%s"
					}
				},
				{
					"term": {
						"auth_id": %d
					}
				}
			]
		}
	}
}`

var AddPhotoUrlRequest = `{
	"doc": {
		"url": "%s"
	}
}`

// search type which indicates if we are searching by tag or by description
type SearchType string

// Photo struct used in elasticsearch.
type PhotoToIndex struct {
	AuthID		uint		`json:"auth_id"`
	BucketID	uint		`json:"bucket_id"`
	ID 			uint		`json:"id"`
	Name 		string		`json:"name"`
	Tags 		[]string	`json:"tags"`
	Url			string		`json:"url"`
	Description string		`json:"description"`
}

// Init elasticsearch client.
func init() {

	host := conf.ServerCfg.Get(constant.ES_HOST)
	port := conf.ServerCfg.Get(constant.ES_PORT)
	esCfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", host, port),
		},
	}
	var err error
	ESClient, err = elasticsearch.NewClient(esCfg)
	if err != nil {
		//log.Fatalln(err)
		utils.AppLogger.Info(err.Error(), zap.String("service", "init()"))
	}
}

// Index a photo in elasticsearch.
func IndexPhoto(photo *Photo) error {

	// the document we want to index
	photoToIndex := PhotoToIndex{
		AuthID: photo.AuthID,
		BucketID: photo.BucketID,
		ID: photo.ID,
		Name: photo.Name,
		Tags: strings.Split(photo.Tag, ";"),
		Url: photo.Url,
		Description: photo.Description,
	}
	body, _ := json.Marshal(&photoToIndex)

	// set up index request
	request := esapi.IndexRequest{
		Index: conf.ServerCfg.Get(constant.ES_PHOTO_INDEX),
		DocumentID: fmt.Sprintf("%d", photoToIndex.ID),
		Body: bytes.NewReader(body),
		Refresh: "true",
	}

	if res, err := request.Do(context.Background(), ESClient); err == nil {
		defer res.Body.Close()
		if res.IsError() {
			//log.Println("Photo indexing error")
			utils.AppLogger.Info(PhotoIndexingError.Error(), zap.String("service", "IndexPhoto()"))
			return PhotoIndexingError
		}
	} else {
		//log.Println(err)
		utils.AppLogger.Info(err.Error(), zap.String("service", "IndexPhoto()"))
		return PhotoIndexingError
	}
	return nil
}

// Add the photo url in elasticsearch.
func AddPhotoUrl(photoID uint, url string) error {
	queryBody := fmt.Sprintf(AddPhotoUrlRequest, url)

	res, err := ESClient.Update(
		conf.ServerCfg.Get(constant.ES_PHOTO_INDEX),
		fmt.Sprintf("%d", photoID),
		strings.NewReader(queryBody),
		)

	if err != nil {
		//log.Println(err)
		utils.AppLogger.Info(err.Error(), zap.String("service", "AddPhotoUrl()"))
		return PhotoUpdateError
	}

	if res.IsError() {
		//log.Println(PhotoUpdateError)
		utils.AppLogger.Info(PhotoUpdateError.Error(), zap.String("service", "AddPhotoUrl()"))
		return PhotoUpdateError
	} else {
		defer res.Body.Close()
		resMap := make(map[string]interface{})
		if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
			//log.Println(err)
			utils.AppLogger.Info(err.Error(), zap.String("service", "AddPhotoUrl()"))
			return PhotoUpdateError
		} else {
			if fmt.Sprintf("%d", resMap["updated"]) != "0" {
				return nil	// update success
			}
		}
	}
	return PhotoUpdateError
}

// Search photo(s) by the given field
// 1. searchType = SEARCH_BY_TAG, the field is a tag
// 2. searchType = SEARCH_BY_DESC, the field is a description
func SearchPhoto(field string, authID uint, offset int, searchType SearchType) ([]PhotoToIndex, error) {
	queryBody := fmt.Sprintf(SearchRequest, searchType, field, authID)
	photos := make([]PhotoToIndex, 0, constant.PAGE_SIZE)

	res, err := ESClient.Search(
		ESClient.Search.WithContext(context.Background()),
		ESClient.Search.WithIndex(conf.ServerCfg.Get(constant.ES_PHOTO_INDEX)),
		ESClient.Search.WithBody(strings.NewReader(queryBody)),
		ESClient.Search.WithFrom(offset),
		ESClient.Search.WithSize(constant.PAGE_SIZE),
		)

	if err != nil {
		//log.Println(err)
		utils.AppLogger.Info(err.Error(), zap.String("service", "SearchPhoto()"))
		return photos, PhotoSearchError
	}

	if res.IsError() {
		//log.Println(PhotoSearchError)
		utils.AppLogger.Info(PhotoSearchError.Error(), zap.String("service", "SearchPhoto()"))
		return photos, PhotoSearchError
	} else {
		defer res.Body.Close()
		resMap := make(map[string]interface{})
		if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
			//log.Println(err)
			utils.AppLogger.Info(err.Error(), zap.String("service", "SearchPhoto()"))
			return photos, PhotoSearchError
		} else {
			// for each hit in the response, we marshal the source into the photo object
			for _, hit := range resMap["hits"].(map[string]interface{})["hits"].([]interface{}) {
				source, _ := json.Marshal(hit.(map[string]interface{})["_source"])
				photo := PhotoToIndex{}
				json.Unmarshal(source, &photo)
				photos = append(photos, photo)
			}
		}
	}

	return photos, nil
}