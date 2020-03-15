/*================================================================
*Copyright (C) 2020 BGBiao Ltd. All rights reserved.
*
*FileName:config.go
*Author:Xuebiao Xu
*Date:2020年03月15日
*Description:
*
================================================================*/
package api

import (
    "fmt"
    "github.com/spf13/viper"
)

type DBConfig struct {
    Host    string    `json:"host"` 
    Port    int64     `json:"port"`
    User    string    `json:"user"`
    Passwd  string    `json:"passwd"`
    Database  string  `json:"database"`
}

var (
    DbConfig DBConfig
)

func ParserConfig() {
    config := viper.New()
    config.AddConfigPath("./config")
    config.SetConfigName("config")
    config.SetConfigType("ini")

    if err := config.ReadInConfig(); err != nil {
        panic(err)
    }

    fmt.Println(config.GetString("db.host"))
    fmt.Println(config.GetInt64("db.port"))
    DbConfig.Host = config.GetString("db.host")
    DbConfig.Port = config.GetInt64("db.port")
    DbConfig.User = config.GetString("db.user")
    DbConfig.Passwd = config.GetString("db.passwd")
    DbConfig.Database = config.GetString("db.database")
}
