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
	ValidMethods         []string // 有效的加密方法列表，如果不为空，则Parse.Method.Alg()必需是VaildMethods的一种，否则报错
	UseJSONNumber        bool     // Use JSON Number format in JSON decoder
	SkipClaimsValidation bool     // 在解析token时跳过claims的验证
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

## 基本用法

### 创建一个token
```go
type MyCustomClaims struct{
	Username string `json:"username"`
    jwt.StandardClaims
}

secretKey := []byte("Helloworld")

// 直接创建一个token对象，加密方式为HS256
// 下面的代码等于
// token :=  NewWithClaims(jwt.SigningMethodHS256,MyCustomClaims{"Mike"})
token := New(jwt.SigningMethodHS256)
claims := MyCustomClaims{
	"Mike",
}
token.Claims = claims

// 获得最终的tokenStr
tokenStr, err := token.SignedString(secretKey)
    
```

### 解析一个token
```go
var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"

token, err := jwt.Parse(tokenString,func(token *jwt.Token)(interface{},error){
	return []byte("Helloworld"), nil
})

// 检查token是否合法
if token.Valid  {
	fmt.Println("token合法")
} else {
	fmt.Println("token不合法, err:",err)
}

```

## 源码分析
首先我们先来看Parse()
```go
func (p *Parser) Parse(tokenString string, keyFunc Keyfunc) (*Token, error) {
    return p.ParseWithClaims(tokenString, MapClaims{}, keyFunc)
}
```
实际上是调用了ParseWithClaims,第二个参数就是一个map[string]interface,这个函数的源码如下,这个函数内部调用的ParseUnverified,我们先来看看这个函数
官方的解释是，这个函数不校验签名的有效性，只单纯负责把tokenStr变成Token对象，而之后的事情就是交给ParseWithClaims来做啦
```go
// WARNING: Don't use this method unless you know what you're doing
//
// This method parses the token but doesn't validate the signature. It's only
// ever useful in cases where you know the signature is valid (because it has
// been checked previously in the stack) and you want to extract values from
// it.
func (p *Parser) ParseUnverified(tokenString string, claims Claims) (token *Token, parts []string, err error) {
	// 将tokenStr的三部分分开,如果不是三部分则报错:不是规范的token
	parts = strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, parts, NewValidationError("token contains an invalid number of segments", ValidationErrorMalformed)
	}
	// 将token的Raw写进去
	token = &Token{Raw: tokenString}

	// parse Header
	// 解析头部
	// 这里的DecodeSegment是base64url编码的译码
	var headerBytes []byte
	if headerBytes, err = DecodeSegment(parts[0]); err != nil {
		// 传进来的token不应该包含bearer
		if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
			return token, parts, NewValidationError("tokenstring should not contain 'bearer '", ValidationErrorMalformed)
		}
		return token, parts, &ValidationError{Inner: err, Errors: ValidationErrorMalformed}
	}
	// 将数据解码到Header中，Header是一个map[string]interface{}
	if err = json.Unmarshal(headerBytes, &token.Header); err != nil {
		return token, parts, &ValidationError{Inner: err, Errors: ValidationErrorMalformed}
	}

	// parse Claims
	// 解析Claims
	var claimBytes []byte
	token.Claims = claims

	// Base64url译码
	if claimBytes, err = DecodeSegment(parts[1]); err != nil {
		return token, parts, &ValidationError{Inner: err, Errors: ValidationErrorMalformed}
	}
	// json解码
	dec := json.NewDecoder(bytes.NewBuffer(claimBytes))
	if p.UseJSONNumber {
		dec.UseNumber()
	}
	// JSON Decode.  Special case for map type to avoid weird pointer behavior
	// 这里传进来的claims是一个自定义结构体的引用,所以decode到claims里面，token.Claims里面也有数据了
	if c, ok := token.Claims.(MapClaims); ok {
		err = dec.Decode(&c)
	} else {
		err = dec.Decode(&claims)
	}
	// Handle decode error
	if err != nil {
		return token, parts, &ValidationError{Inner: err, Errors: ValidationErrorMalformed}
	}

	// Lookup signature method
	// 从Header里面取出alg,获得SigningMethod
	if method, ok := token.Header["alg"].(string); ok {
		if token.Method = GetSigningMethod(method); token.Method == nil {
			return token, parts, NewValidationError("signing method (alg) is unavailable.", ValidationErrorUnverifiable)
		}
	} else {
		return token, parts, NewValidationError("signing method (alg) is unspecified.", ValidationErrorUnverifiable)
	}

	return token, parts, nil
}
```
可以看到，ParseUnverified这个方法真的只是单纯的解码Header段和Claim段，然后检查一下用的alg是不是合法，就返回了，让我们继续往下看验证的逻辑

```go
func (p *Parser) ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {
	
	// 这里我们得到了未验证的token对象
	token, parts, err := p.ParseUnverified(tokenString, claims)
	if err != nil {
		return token, err
	}

	
	// Verify signing method is in the required set
	// 验证SigningMethod，如果提供了VaildMethods字段，只有这几种SigningMethod才被认为是合法的
	if p.ValidMethods != nil {
		var signingMethodValid = false
		var alg = token.Method.Alg()
		for _, m := range p.ValidMethods {
			if m == alg {
				signingMethodValid = true
				break
			}
		}
		if !signingMethodValid {
			// signing method is not in the listed set
			return token, NewValidationError(fmt.Sprintf("signing method %v is invalid", alg), ValidationErrorSignatureInvalid)
		}
	}

	// Lookup key
	// 检查KeyFunc,这里它回去检查KeyFunc(token)的err，可是我看网上的教程都是直接返回SecretKey, nil的
	// 我不是很理解这段代码的意思
	var key interface{}
	if keyFunc == nil {
		// keyFunc was not provided.  short circuiting validation
		return token, NewValidationError("no Keyfunc was provided.", ValidationErrorUnverifiable)
	}
	if key, err = keyFunc(token); err != nil {
		// keyFunc returned an error
		if ve, ok := err.(*ValidationError); ok {
			return token, ve
		}
		return token, &ValidationError{Inner: err, Errors: ValidationErrorUnverifiable}
	}

	vErr := &ValidationError{}

	// Validate Claims
	// 验证Claims段,如果SkipClaimsValidation为false
	if !p.SkipClaimsValidation {
		// Claims是一个实现了Vaild方法的接口,所以直接调用token.Claims的Valid方法
		// jwt.Standardclaims实现了Claims结构
		if err := token.Claims.Valid(); err != nil {
            
			// If the Claims Valid returned an error, check if it is a validation error,
			// If it was another error type, create a ValidationError with a generic ClaimsInvalid flag set
			if e, ok := err.(*ValidationError); !ok {
				vErr = &ValidationError{Inner: err, Errors: ValidationErrorClaimsInvalid}
			} else {
				vErr = e
			}
		}
	}

	// Perform validation
	// 这里是调用SigningMethod结构的Verify方法，检查token是否有效
	token.Signature = parts[2]
	if err = token.Method.Verify(strings.Join(parts[0:2], "."), token.Signature, key); err != nil {
		vErr.Inner = err
		vErr.Errors |= ValidationErrorSignatureInvalid
	}
    
	// vaild就是看vErr.Error是否为0，这里面记录着验证token时可能出现的一些错误
	if vErr.valid() {
		token.Valid = true
		return token, nil
	}

	return token, vErr
}

```
ok,关于解析token的主要方法我们已经看完了，接下来我们来看看如何生成一个token，其实就是反着操作一遍
先看New函数，选择一种SigningMethod，新建一个token，内部调用NewWithClaims
```go
// Create a new Token.  Takes a signing method
func New(method SigningMethod) *Token {
	return NewWithClaims(method, MapClaims{})
}
```
再看NewWithClaims,发现就是简单的给JwtToken的三个部分赋值
```go
func NewWithClaims(method SigningMethod, claims Claims) *Token {
	return &Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
		},
		Claims: claims,
		Method: method,
	}
}
```
最后是SignedString，即使用alg的算法给token加密，生成最终的tokenStr,内部调用了SigningString，所以先看SigningString
发现SigningString就是把token的头部先变成json然后base64url编码，但是没有生成jwtToken的最后一个部分
```go
// Generate the signing string.  This is the
// most expensive part of the whole deal.  Unless you
// need this for something special, just go straight for
// the SignedString.
func (t *Token) SigningString() (string, error) {
    var err error
    parts := make([]string, 2)
    for i, _ := range parts {
        var jsonValue []byte
        if i == 0 {
            if jsonValue, err = json.Marshal(t.Header); err != nil {
                return "", err
            }
        } else {
            if jsonValue, err = json.Marshal(t.Claims); err != nil {
                return "", err
            }
        }

    parts[i] = EncodeSegment(jsonValue)
    }
    return strings.Join(parts, "."), nil
}
```
所以SignedString作用就是用给定的加密方法和你的SecretKey对前面两部分加密，添在token的最后一段，至此token生成完毕
```go
// Get the complete, signed token
func (t *Token) SignedString(key interface{}) (string, error) {
	var sig, sstr string
	var err error
	if sstr, err = t.SigningString(); err != nil {
		return "", err
	}
	if sig, err = t.Method.Sign(sstr, key); err != nil {
		return "", err
	}
	return strings.Join([]string{sstr, sig}, "."), nil
}

```

