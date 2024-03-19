package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SqlLimitObject struct {

	// SQL限流ID。
	Id string `json:"id"`

	// 由SQL的语法解析树计算出的内部哈希码，默认为0，取值范围（-9223372036854775808~ 9223372036854775807）。
	QueryId string `json:"query_id"`

	// SQL语句的文本形式。
	QueryString string `json:"query_string"`

	// 同时执行的SQL数量，小于等于0表示不限制，默认为0，取值范围（-1~50000）。
	MaxConcurrency int32 `json:"max_concurrency"`

	// 是否生效
	IsEffective bool `json:"is_effective"`

	// 最大等待时间，单位为秒。
	MaxWaiting int32 `json:"max_waiting"`

	// 为不是模式限定的名称设置模式搜索顺序，默认为public。
	SearchPath string `json:"search_path"`
}

func (o SqlLimitObject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SqlLimitObject struct{}"
	}

	return strings.Join([]string{"SqlLimitObject", string(data)}, " ")
}
