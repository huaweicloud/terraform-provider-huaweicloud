package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AddressAndPorts struct {

	// IP地址或域名。
	Address string `json:"address"`

	// 端口号。
	Port *int32 `json:"port,omitempty"`
}

func (o AddressAndPorts) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddressAndPorts struct{}"
	}

	return strings.Join([]string{"AddressAndPorts", string(data)}, " ")
}
