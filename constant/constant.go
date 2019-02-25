package constant

const (

	// JWT constants
	JWT_SECRET 			= "JWT_SECRET"
	JWT 				= "jwt"
	JWT_EXP_MINUTE 		= 30
	PHOTO_STORAGE_ADMIN = "admin"

	// Server constants
	SERVER_PORT = "SERVER_PORT"
	SERVER_DOMAIN = "SERVER_DOMAIN"
	SERVER_PATH = "SERVER_PATH"
	PAGE_SIZE 	= 20

	// DB constants
	DB_CONNECT 	= "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
	DB_TYPE 	= "DB_TYPE"
	DB_HOST 	= "DB_HOST"
	DB_PORT 	= "DB_PORT"
	DB_USER 	= "DB_USER"
	DB_PWD 		= "DB_PWD"
	DB_NAME 	= "DB_NAME"

	// Redis constants
	REDIS_HOST = "REDIS_HOST"
	REDIS_PORT = "REDIS_PORT"

	// Auth constants
	COOKIE_MAX_AGE 	= 1800
	LOGIN_MAX_AGE 	= 1800
	LOGIN_USER 		= "LOGIN_"

	// COS constants
	COS_BUCKET_NAME = "COS_BUCKET"
	COS_APP_ID 		= "COS_APP_ID"
	COS_REGION 		= "COS_REGION"
	COS_SECRET_ID 	= "COS_SECRET_ID"
	COS_SECRET_KEY 	= "COS_SECRET_KEY"

	// Callback constants
	URL_UPDATE_CHANNEL 		= "PHOTO_URL_UPDATE"
	PHOTO_UPDATE_ID_FORMAT 	= "photo-%d"
	PHOTO_DELETE_CHANNEL 	= "PHOTO_DELETE"
)
