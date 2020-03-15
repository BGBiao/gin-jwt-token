/*
Copyright 2019 BGBiao Ltd. All rights reserved.
@File   : rds-disk.go
@Time   : 2019/11/11 18:22:56
@Update : 2019/11/11 18:22:56
@Author : BGBiao
@Version: 1.0
@Contact: weichaungxxb@qq.com
@Desc   : None
*/

package controller

import (
	_ "fmt"
	"log"
	"net/http"
	"time"
	md "warnning-trigger/middleware"
	"warnning-trigger/model"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 用户注册信息
// 注意:注册信息可以使用Gin内部的校验工具或者beego的校验工具进行校验
type RegisterInfo struct {
	Phone int64  `json:"phone"`
	Name  string `json:"name"`
	Pwd   string `json:"password"`
	Email string `json:"email"`
}

// 用户注册接口
func RegisterUser(c *gin.Context) {
	var registerInfo RegisterInfo
	bindErr := c.BindJSON(&registerInfo)
	if bindErr == nil {
		// 用户注册
		err := model.Register(registerInfo.Name, registerInfo.Pwd, registerInfo.Phone, registerInfo.Email)

		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 0,
				"msg":    "success ",
				"data":   nil,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "注册失败" + err.Error(),
				"data":   nil,
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "用户注册解析数据失败" + bindErr.Error(),
			"data":   nil,
		})
	}
}

// 登陆结果
type LoginResult struct {
	Token string `json:"token"`
	// 用户模型
	Name string `json:"name"`
	//model.User
}

// 登陆接口 用户名和密码登陆
// name,password
func Login(c *gin.Context) {
	var loginReq model.LoginReq
	if c.BindJSON(&loginReq) == nil {
		// 登陆逻辑校验
		isPass, user, err := model.LoginCheck(loginReq)
		if isPass {
			generateToken(c, user)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "验证失败" + err.Error(),
				"data":   nil,
			})
		}

	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "用户数据解析失败",
			"data":   nil,
		})
	}
}

// token生成器
func generateToken(c *gin.Context, user model.User) {
	// 构造SignKey: 签名和解签名需要使用一个值
	j := md.NewJWT()

	// 构造用户claims信息(负荷)
	claims := md.CustomClaims{
		user.Name,
		user.Email,
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 签名过期时间
			Issuer:    "bgbiao.top",                    // 签名颁发者
		},
	}

	// 根据claims生成token对象
	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
			"data":   nil,
		})
	}

	log.Println(token)
	// 获取用户相关数据
	data := LoginResult{
		Name:  user.Name,
		Token: token,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "登陆成功",
		"data":   data,
	})
	return
}

// 测试一个需要认证的接口
func GetDataByTime(c *gin.Context) {
	claims := c.MustGet("claims").(*md.CustomClaims)
	if claims != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "token有效",
			"data":   claims,
		})
	}
}
