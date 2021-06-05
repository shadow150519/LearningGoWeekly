# 使用net/http发送request

## Get方法
```go
// 普通的Get
resp, err := http.Get("https://www.baidu.com")

// 带url参数
name := "wtx"
age := "18"
params := url.Values{}
Url, _ := url.Parse("https://baidu.com")
params.Set("name",name)
params.Set("age",age)
Url.RawQuery = params.Encode()
resp, err = http.Get(Url.String())

// 进行请求头等更进一步的配置
client := &http.Client{}
req,_ := http.NewRequest("GET","https://baidu.com",nil)
req.Header.Add("name","zhaofan")
req.Header.Add("age","3")
resp,_ := client.Do(req)
```

## Post方法
```go
resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
...
resp, err := http.PostForm("http://example.com/form",
url.Values{"key": {"Value"}, "id": {"123"}})
```