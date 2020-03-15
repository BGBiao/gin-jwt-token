/*================================================================
*Copyright (C) 2020 BGBiao Ltd. All rights reserved.
*
*FileName:user.go
*Author:Xuebiao Xu
*Date:2020年03月10日
*Description:
*
================================================================*/
package model

import (
	"fmt"
	"time"

)

// 构造用户表
type User struct {
	Id          int32  `gorm:"AUTO_INCREMENT"`
	Name        string `json:"name"`
	Pwd         string `json:"password"`
	Phone       int64  `gorm:"DEFAULT:0"`
	Email       string `gorm:"type:varchar(20);unique_index;"`
	CreatedAt *time.Time
	UpdateTAt  *time.Time
}

// LoginReq请求参数
type LoginReq struct {
	Name string `json:"name"`
	Pwd  string `json:"password"`
}

// 初始化表结构以及定义

func InitModel() {
    DB.AutoMigrate(&User{})
}

// 插入数据
func (user *User) Insert() error {
	//这里使用了Table()函数，如果你没有指定全局表名禁用复数，或者是表名跟结构体名不一样的时候
	//你可以自己在sql中指定表名。这里是示例，本例中这个函数可以去除。
	// 需要注意的是Create函数的参数必须是指针
	// return DB.Table("user").Create(user).Error
	return DB.Model(&User{}).Create(&user).Error
}


// 用户注册
func Register(username, pwd string, phone int64, email string) error {
	fmt.Println(username, pwd, phone, email)

	if CheckUser(username) {
		return fmt.Errorf("用户已经存在，请直接登陆")
	}

	// defer db.Close()
	// 需要生成一个uuid: Id为自增
	// 构造用户注册信息
	user := User{
		Name:  username,
		Pwd:   pwd,
		Phone: phone,
		Email: email,
	}
	insertErr := user.Insert()
	return insertErr

}

// 用户检查
func CheckUser(username string) bool {

	result := false
	// 指定库
	var user User

	dbResult := DB.Where("name = ?", username).Find(&user)
	if dbResult.Error != nil {
		fmt.Printf("获取用户信息失败:%v\n", dbResult.Error)
	} else {
		result = true
	}
	return result
}

// LoginCheck验证
func LoginCheck(login LoginReq) (bool, User, error) {
	userData := User{}
	userExist := false

	var user User
	dbErr := DB.Where("name = ?", login.Name).Find(&user).Error

	if dbErr != nil {
		return userExist, userData, dbErr
	}
	if login.Name == user.Name && login.Pwd == user.Pwd {
		userExist = true
		userData.Name = user.Name
		userData.Email = user.Email
	}

	if !userExist {
		return userExist, userData, fmt.Errorf("登陆信息有误")
	}
	return userExist, userData, nil
}


