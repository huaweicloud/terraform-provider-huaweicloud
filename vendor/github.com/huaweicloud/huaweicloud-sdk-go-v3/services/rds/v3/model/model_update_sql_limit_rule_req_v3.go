package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateSqlLimitRuleReqV3 struct {

	// 数据库名称。
	DbName string `json:"db_name"`

	// SQL限流ID。
	Id string `json:"id"`

	// 同时执行的sql数量，小于等于0表示不限制，默认为0，取值范围（-1~50000）。
	MaxConcurrency int32 `json:"max_concurrency"`

	// 最大等待时间，单位为秒。
	MaxWaiting int32 `json:"max_waiting"`
}

func (o UpdateSqlLimitRuleReqV3) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSqlLimitRuleReqV3 struct{}"
	}

	return strings.Join([]string{"UpdateSqlLimitRuleReqV3", string(data)}, " ")
}
