package constant

const (

	// User related responses
	USER_ALREADY_EXIST 		= 1001
	USER_ADD_SUCCESS 		= 1002
	USER_AUTH_SUCCESS 		= 1003
	USER_AUTH_ERROR 		= 1004
	USER_AUTH_TIMEOUT 		= 1005
	USER_SIGNOUT_SUCCESS 	= 1006

	// JWT related responses
	JWT_GENERATION_ERROR 	= 2001
	JWT_MISSING_ERROR 		= 2002
	JWT_PARSE_ERROR 		= 2003

	// Bucket related responses
	BUCKET_ALREADY_EXIST 	= 3001
	BUCKET_ADD_SUCCESS 		= 3002
	BUCKET_NOT_EXIST 		= 3003
	BUCKET_DELETE_SUCCESS 	= 3004
	BUCKET_UPDATE_SUCCESS 	= 3005
	BUCKET_GET_SUCCESS 		= 3006

	// Photo related responses
	PHOTO_ALREADY_EXIST 			= 4001
	PHOTO_ADD_IN_PROCESS 			= 4002
	PHOTO_UPLOAD_SUCCESS 			= 4003
	PHOTO_UPLOAD_ERROR 				= 4004
	PHOTO_NOT_EXIST 				= 4005
	PHOTO_DELETE_SUCCESS 			= 4006
	PHOTO_UPDATE_SUCCESS 			= 4007
	PHOTO_GET_SUCCESS 				= 4008
	PHOTO_SEARCH_BY_TAG_SUCCESS 	= 4009
	PHOTO_SEARCH_BY_DESC_SUCCESS	= 4010

	// Internal server responses
	INTERNAL_SERVER_ERROR 	= 5001
	PAGINATION_SUCCESS 		= 8001
	INVALID_PARAMS 			= 9001

)

var Message map[int]string

// Init the message map.
func init() {
	Message = make(map[int]string)
	Message[INVALID_PARAMS] 		= "Invalid parameters."
	Message[USER_ALREADY_EXIST] 	= "User already exists."
	Message[USER_ADD_SUCCESS] 		= "Add user success."
	Message[USER_AUTH_SUCCESS] 		= "User authentication success."
	Message[USER_AUTH_ERROR] 		= "User authentication fail."
	Message[USER_AUTH_TIMEOUT] 		= "User authentication timeout."
	Message[USER_SIGNOUT_SUCCESS] 	= "User sign out success."
	Message[JWT_GENERATION_ERROR] 	= "JWT generation fail."
	Message[JWT_MISSING_ERROR] 		= "JWT is missing."
	Message[JWT_PARSE_ERROR]		= "JWT parse error."
	Message[INTERNAL_SERVER_ERROR] 	= "Internal server error."
	Message[BUCKET_ALREADY_EXIST] 	= "Bucket already exists."
	Message[BUCKET_ADD_SUCCESS] 	= "Add bucket success."
	Message[BUCKET_NOT_EXIST]		= "Bucket does not exist."
	Message[BUCKET_DELETE_SUCCESS] 	= "Bucket delete success."
	Message[BUCKET_UPDATE_SUCCESS] 	= "Bucket update success."
	Message[BUCKET_GET_SUCCESS] 	= "Bucket get success."
	Message[PHOTO_ALREADY_EXIST] 	= "Photo already exists."
	Message[PHOTO_ADD_IN_PROCESS] 	= "Adding photo is in process."
	Message[PHOTO_UPLOAD_SUCCESS] 	= "Photo upload success."
	Message[PHOTO_UPLOAD_ERROR] 	= "Photo upload error."
	Message[PHOTO_NOT_EXIST] 		= "Photo does not exist."
	Message[PHOTO_DELETE_SUCCESS] 	= "Photo delete success."
	Message[PHOTO_UPDATE_SUCCESS]	= "Photo update success."
	Message[PHOTO_GET_SUCCESS]		= "Photo get success."
	Message[PHOTO_SEARCH_BY_TAG_SUCCESS] = "Photo search by tag success."
	Message[PHOTO_SEARCH_BY_DESC_SUCCESS] = "Photo search by description success."
}

// Translate a response code to a detailed message.
func GetMessage(code int) string {
	msg, ok := Message[code]
	if ok {
		return msg
	}
	return ""
}