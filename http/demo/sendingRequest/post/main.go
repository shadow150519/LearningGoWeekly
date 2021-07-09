package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	body := struct {
		Username string `json:"username"`
		Age int `json:"age"`
	}{"wtx",12}
	bodyBytes, _ := json.Marshal(&body)
	resp, err :=http.Post("https://baidu.com","application/json",bytes.NewReader(bodyBytes))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.StatusCode)
	var respBuffer []byte
	resp.Body.Read(respBuffer)
	fmt.Println(string(respBuffer))

}
