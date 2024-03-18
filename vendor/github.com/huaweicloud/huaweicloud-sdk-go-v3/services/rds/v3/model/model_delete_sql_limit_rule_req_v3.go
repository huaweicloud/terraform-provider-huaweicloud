package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeleteSqlLimitRuleReqV3 struct {

	// 数据库名称。
	DbName string `json:"db_name"`

	// SQL限流ID。
	Id string `json:"id"`
}

func (o DeleteSqlLimitRuleReqV3) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSqlLimitRuleReqV3 struct{}"
	}

	return strings.Join([]string{"DeleteSqlLimitRuleReqV3", string(data)}, " ")
}
