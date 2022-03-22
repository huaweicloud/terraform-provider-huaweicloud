package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlProxyAvailable struct {
	// 可用区编码。

	Code *string `json:"code,omitempty"`
	// 可用区描述。

	Description *string `json:"description,omitempty"`
}

func (o MysqlProxyAvailable) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlProxyAvailable struct{}"
	}

	return strings.Join([]string{"MysqlProxyAvailable", string(data)}, " ")
}
