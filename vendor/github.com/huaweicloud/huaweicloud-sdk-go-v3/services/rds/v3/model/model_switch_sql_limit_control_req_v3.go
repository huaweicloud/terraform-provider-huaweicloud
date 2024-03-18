package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SwitchSqlLimitControlReqV3 struct {

	// 数据库名称。
	DbName string `json:"db_name"`

	// SQL限流ID。
	Id string `json:"id"`

	// SQL限流动作标志。 取值为“open”：表示开启当前SQL限流。 取值为“close”：表示关闭当前SQL限流。 取值为“disable_all”：表示禁用所有SQL限流。
	Action string `json:"action"`
}

func (o SwitchSqlLimitControlReqV3) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SwitchSqlLimitControlReqV3 struct{}"
	}

	return strings.Join([]string{"SwitchSqlLimitControlReqV3", string(data)}, " ")
}
