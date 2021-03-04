package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	Uname string `json:"username"`
	jwt.StandardClaims
}

func main() {
	sec := []byte("123abc")
	token_obj := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		Uname:          "shenyi",
		StandardClaims: jwt.StandardClaims{},
	})
	token, _ := token_obj.SignedString(sec)
	fmt.Println(token)

	uc := UserClaim{}


	getToken, _ := jwt.ParseWithClaims(token,&uc, func(token *jwt.Token) (i interface{}, err error) {
		return sec, nil
	})
	if getToken.Valid {
		fmt.Println(getToken.Claims.(*UserClaim).Uname)
		fmt.Println(getToken.Claims.(*UserClaim).ExpiresAt)
	}
}
//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InNoZW55aSJ9.5tRNAqDj3gK5Zm8TPTyXGyCcxKWtpnhKdVcfZOo5qh4
