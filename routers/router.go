package routers

import (
	"gin-photo-storage/apis/v1"
	"gin-photo-storage/middleware"
	"github.com/gin-gonic/gin"
	_ "gin-photo-storage/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// a global router
var Router *gin.Engine

// Init router, adding paths to it.
func init() {
	Router = gin.Default()
	checkAuthMdw := middleware.GetAuthMiddleware()			// middleware for authentication
	refreshMdw := middleware.GetRefreshMiddleware()			// middleware for refresh auth token
	paginationMdw := middleware.GetPaginationMiddleware()	// middleware for pagination

	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// api group for v1
	v1Group := Router.Group("/api/v1")
	{
		// api group for authentication
		authGroup := v1Group.Group("/auth")
		{
			authGroup.POST("/add", v1.AddAuth)
			authGroup.POST("/check", v1.CheckAuth)
		}

		// api group for bucket
		bucketGroup := v1Group.Group("/bucket")
		{
			// must check auth & refresh auth token before any operation
			bucketGroup.POST("/add", checkAuthMdw, refreshMdw, v1.AddBucket)
			bucketGroup.DELETE("/delete", checkAuthMdw, refreshMdw, v1.DeleteBucket)
			bucketGroup.PUT("/update", checkAuthMdw, refreshMdw, v1.UpdateBucket)
			bucketGroup.GET("/get_by_id", checkAuthMdw, refreshMdw, v1.GetBucketByID)
			bucketGroup.GET("/get_by_auth_id", checkAuthMdw, refreshMdw, paginationMdw, v1.GetBucketByAuthID)
		}

		// api group for photo
		photoGroup := v1Group.Group("/photo")
		{
			// must check auth & refresh auth token before any operation
			photoGroup.POST("/add", checkAuthMdw, refreshMdw, v1.AddPhoto)
			photoGroup.GET("/upload_status", checkAuthMdw, refreshMdw, v1.GetPhotoUploadStatus)
			photoGroup.DELETE("/delete", checkAuthMdw, refreshMdw, v1.DeletePhoto)
			photoGroup.PUT("/update", checkAuthMdw, refreshMdw, v1.UpdatePhoto)
			photoGroup.GET("/get_by_id", checkAuthMdw, refreshMdw, v1.GetPhotoByID)
			photoGroup.GET("/get_by_bucket_id", checkAuthMdw, refreshMdw, paginationMdw, v1.GetPhotoByBucketID)
			photoGroup.GET("/search", checkAuthMdw, refreshMdw, paginationMdw, v1.SearchPhoto)
		}
	}
}
