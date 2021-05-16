# jwt-go笔记

## jwt介绍
jwt(json web token)是一种用于前后端身份认证的方法,一个jwt由header，payload，和signature组成。

+ header:包含了token类型和算法类型
+ payload:包含了一些用户自定义或jwt预定义的一些数据，每一个数据叫一个claim,**注意不要将敏感信息放入**
+ signature:将header和payload经过base64编码,加上一个secret密钥，整体经过header中的算法加密后生成

## jwt-go重要的几个结构

1.Claims 
```go
type Claims interface {
	Valid() error
}
```
claims是一个实现了Valid方法的interface，Valid方法用于判断该claim是否合法

2.Keyfunc
```go
type Keyfunc func(*Token) (interface{}, error)
```
Keyfunc在使用时一般都是返回secret密钥，可以根据Token的种类不同返回不同的密钥.

`官方文档:This allows you to use properties in the Header of the token (such as 'kid') to identify which key to use.`

3.Mapclaims
```go
type MapClaims map[string]interface{}
```
一个用于放decode出来的claim的map,有Vaild和一系列VerifyXXX的方法

4.Parser
```go
type Parser struct {
	ValidMethods         []string // If populated, only these methods will be considered valid
	UseJSONNumber        bool     // Use JSON Number format in JSON decoder
	SkipClaimsValidation bool     // Skip claims validation during token parsing
}
```
用来将tokenstr转换成token


5.SigningMethod
```go
type SigningMethod interface {
	Verify(signingString, signature string, key interface{}) error // Returns nil if signature is valid
	Sign(signingString string, key interface{}) (string, error)    // Returns encoded signature or error
	Alg() string                                                   // returns the alg identifier for this method (example: 'HS256')
}
```
签名方法的接口，可以通过实现这个接口自定义签名方法，jwt-go内置一些实现了SigningMethod的结构体

6.StandardClaims
```go
type StandardClaims struct {
	Audience  string `json:"aud,omitempty"` 
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}
```
jwt官方规定的一些预定义的payload:
+ iss (issuer)：签发人
+ exp (expiration time)：过期时间
+ sub (subject)：主题
+ aud (audience)：受众
+ nbf (Not Before)：生效时间
+ iat (Issued At)：签发时间
+ jti (JWT ID)：编号

7.Token
```go
type Token struct {
	Raw       string                 // The raw token.  Populated when you Parse a token
	Method    SigningMethod          // The signing method used or to be used
	Header    map[string]interface{} // The first segment of the token
	Claims    Claims                 // The second segment of the token
	Signature string                 // The third segment of the token.  Populated when you Parse a token
	Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
}
```
Token的结构体

8.ValidationError
```go
type ValidationError struct {
	Inner  error  // stores the error returned by external dependencies, i.e.: KeyFunc
	Errors uint32 // bitfield.  see ValidationError... constants
	// contains filtered or unexported fields
}
```
定义解析Token时遇到的一些错误