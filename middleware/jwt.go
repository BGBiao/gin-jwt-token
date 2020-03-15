package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "bgbiao.top" // 签名信息应该设置成动态从库中获取
)

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "请求未携带token，无权限访问",
				"data":   nil,
			})
			c.Abort()
			return
		}

		log.Print("get token: ", token)
		j := NewJWT()
		// 解析token中包含的相关信息
		claims, err := j.ParserToken(token)

		fmt.Println(claims, err)
		if err != nil {
			// token过期
			if err == TokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    "token授权已过期，请重新申请授权",
					"data":   nil,
				})
				c.Abort()
				return
			}
			// 其他错误
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
				"data":   nil,
			})
			c.Abort()
			return
		}

		// 解析到具体的claims相关信息
		c.Set("claims", claims)

	}
}

// JWT基本数据结构
// 签名的signkey
type JWT struct {
	SigningKey []byte
}

// 定义载荷
type CustomClaims struct {
	Name  string `json:"userName"`
	Email string `json:"email"`
	// StandardClaims结构体实现了Claims接口(Valid()函数)
	jwt.StandardClaims
}

// 初始化JWT实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// 获取signkey(这里写死成一个变量了)
func GetSignKey() string {
	return SignKey
}

func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

// 创建Token(基于用户的基本信息claims)
// 使用HS256算法进行token生成
// 使用用户基本信息claims以及签名key(signkey)生成token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#Token
	// 返回一个token的结构体指针
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// token解析
// Couldn't handle this token:
func (j *JWT) ParserToken(tokenString string) (*CustomClaims, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#ParseWithClaims
	// 输入用户自定义的Claims结构体对象,token,以及自定义函数来解析token字符串为jwt的Token结构体指针
	// Keyfunc是匿名函数类型: type Keyfunc func(*Token) (interface{}, error)
	// func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	fmt.Println(token, err)
	if err != nil {
		// https://gowalker.org/github.com/dgrijalva/jwt-go#ValidationError
		// jwt.ValidationError 是一个无效token的错误结构
		if ve, ok := err.(*jwt.ValidationError); ok {
			// ValidationErrorMalformed是一个uint常量，表示token不可用
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
				// ValidationErrorExpired表示Token过期
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
				// ValidationErrorNotValidYet表示无效token
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}

		}
	}

	// 将token中的claims信息解析出来和用户原始数据进行校验
	// 做以下类型断言，将token.Claims转换成具体用户自定义的Claims结构体
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid

}

// 更新Token
func (j *JWT) UpdateToken(tokenString string) (string, error) {
	// TimeFunc为一个默认值是time.Now的当前时间变量,用来解析token后进行过期时间验证
	// 可以使用其他的时间值来覆盖
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	// 拿到token基础数据
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil

	})

	// 校验token当前还有效
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		// 修改Claims的过期时间(int64)
		// https://gowalker.org/github.com/dgrijalva/jwt-go#StandardClaims
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", fmt.Errorf("token获取失败:%v", err)
}
