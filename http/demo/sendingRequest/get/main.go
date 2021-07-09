package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	// 简单的Get
	resp, err := http.Get("https://www.baidu.com")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.StatusCode)
	byteBuff, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(byteBuff))

	// 通过使用变量来生成url query参数来发送Get请求
	params := url.Values{}
	Url, _ := url.Parse("https://baidu.com")
	params.Set("name","wtx")
	params.Set("age","20")
	Url.RawQuery = params.Encode()
	resp, err = http.Get(Url.String())

}
