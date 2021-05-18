package main

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	Name string
	jwt.StandardClaims
}

func main()  {
 	secretKey := []byte("hello")
	str, err := jwt.NewWithClaims(jwt.SigningMethodHS256,&User{"aaa",jwt.StandardClaims{}}).SignedString(secretKey)
	if err != nil {
		fmt.Println("创建token错误",err)
		return
	}
	token, err := jwt.ParseWithClaims(str,&User{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	if user, ok := token.Claims.(*User); ok {
		fmt.Println(user.Name)
	}
}
