package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PrivateIp 服务器私有IP
type PrivateIp struct {
}

func (o PrivateIp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PrivateIp struct{}"
	}

	return strings.Join([]string{"PrivateIp", string(data)}, " ")
}
