package utils

import (
	"gin-photo-storage/conf"
	"gin-photo-storage/constant"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// self-defined user claim
type UserClaim struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

// Generate a JWT based on the user name.
func GenerateJWT(userName string) (string, error) {
	// define a user claim
	claim := UserClaim{
		userName,
		jwt.StandardClaims{
			Issuer:    constant.PHOTO_STORAGE_ADMIN,
			ExpiresAt: time.Now().Add(constant.JWT_EXP_MINUTE * time.Minute).Unix(),
		},
	}

	// generate the claim and the digital signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtString, err := token.SignedString([]byte(conf.ServerCfg.Get(constant.JWT_SECRET)))
	if err != nil {
		log.Fatalln("JWT generation error")
		return "", err
	}
	return jwtString, nil
}

// Parse a JWT into a user claim.
func ParseJWT(jwtString string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(jwtString, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.ServerCfg.Get(constant.JWT_SECRET)), nil
	})

	if token != nil && err == nil {
		if claim, ok := token.Claims.(*UserClaim); ok && token.Valid {
			return claim, nil
		}
	}
	return nil, err
}