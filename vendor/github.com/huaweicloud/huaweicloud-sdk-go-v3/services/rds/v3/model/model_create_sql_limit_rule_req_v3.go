package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateSqlLimitRuleReqV3 struct {

	// 数据库名称。
	DbName string `json:"db_name"`

	// 由SQL的语法解析树计算出的内部哈希码，默认为0，取值范围（-9223372036854775808~ 9223372036854775807）。
	QueryId *int64 `json:"query_id,omitempty"`

	// SQL语句的文本形式。（注：query_id与query_string只可以存在一个）
	QueryString *string `json:"query_string,omitempty"`

	// 同时执行的SQL数量，小于等于0表示不限制，默认为0，取值范围（-1~50000）。
	MaxConcurrency int32 `json:"max_concurrency"`

	// 最大等待时间，单位为秒。
	MaxWaiting int32 `json:"max_waiting"`

	// 为不是模式限定的名称设置模式搜索顺序，默认为public。
	SearchPath *string `json:"search_path,omitempty"`
}

func (o CreateSqlLimitRuleReqV3) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSqlLimitRuleReqV3 struct{}"
	}

	return strings.Join([]string{"CreateSqlLimitRuleReqV3", string(data)}, " ")
}
