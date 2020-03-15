/*================================================================
*Copyright (C) 2020 BGBiao Ltd. All rights reserved.
*
*FileName:config_test.go
*Author:Xuebiao Xu
*Date:2020年03月15日
*Description:
*
================================================================*/
package api

import (
    "fmt"
    "testing"
)

func TestParserConfig(t *testing.T) {
    ParserConfig()  
    if DbConfig.Host == "127.0.0.1" {
        fmt.Println("测试通过")
    }
}
