package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 设备搜索请求结构体。
type SearchSql struct {

	// 搜索sql语句，具体使用方法见类SQL语法使用说明章节
	Sql string `json:"sql"`
}

func (o SearchSql) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SearchSql struct{}"
	}

	return strings.Join([]string{"SearchSql", string(data)}, " ")
}
