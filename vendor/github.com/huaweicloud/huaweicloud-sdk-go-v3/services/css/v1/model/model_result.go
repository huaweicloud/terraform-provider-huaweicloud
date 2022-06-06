package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Result struct {

	// IP地址或域名。
	Address *string `json:"address,omitempty"`

	// 端口号。
	Port *int32 `json:"port,omitempty"`

	// 测试结果。 - 1：表示连接成功。 - 0：表示地址不可达。 - 2：表示端口不可达。 - 3：表示域名无法解析。 - -2：表示位置错误。
	Status *int32 `json:"status,omitempty"`
}

func (o Result) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Result struct{}"
	}

	return strings.Join([]string{"Result", string(data)}, " ")
}
