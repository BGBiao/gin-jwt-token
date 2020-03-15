/*================================================================
*Copyright (C) 2020 BGBiao Ltd. All rights reserved.
*
*FileName:utils_test.go
*Author:Xuebiao Xu
*Date:2020年03月15日
*Description:
*
================================================================*/
package model

import (
    "fmt"
    "testing"
)

func TestInitDb(t *testing.T) {
    connErr := InitMySQLCon()
    if connErr == nil {
        fmt.Println("数据库连接测试通过")
    }else {
        fmt.Printf("数据库练级异常:%v",connErr)
    }

    DB.Close()
}

