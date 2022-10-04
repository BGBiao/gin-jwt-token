/*
Copyright 2019 BGBiao Ltd. All rights reserved.
@File   : main.go
@Time   : 2019/11/11 16:20:45
@Update : 2019/11/11 16:20:45
@Author : BGBiao
@Version: 1.0
@Contact: CloudNativeOps
@Desc   : None
*/
package main

import (
	"expvar"
	"gin-jwt-token/controller"
	md "gin-jwt-token/middleware"
	"gin-jwt-token/model"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化db
	dbErr := model.InitMySQLCon()
	if dbErr != nil {
		panic(dbErr)
	}

	model.InitModel()
	defer model.DB.Close()

	// 初始化一个gin 实例
	router := gin.Default()
	// 使用expvar 暴露内部指标
	router.GET("/debug/vars", gin.WrapH(expvar.Handler()))

	// 定义api
	v1 := router.Group("/apis/v1/")
	{
		v1.POST("/register", controller.RegisterUser)
		v1.POST("/login", controller.Login)
	}

	// secure v1
	sv1 := router.Group("/apis/v1/auth/")
	sv1.Use(md.JWTAuth())
	{
		sv1.GET("/time", controller.GetDataByTime)

	}
	router.Run(":8081")
}
