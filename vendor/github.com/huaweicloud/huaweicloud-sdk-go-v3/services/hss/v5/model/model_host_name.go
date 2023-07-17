package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HostName 服务器名称
type HostName struct {
}

func (o HostName) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HostName struct{}"
	}

	return strings.Join([]string{"HostName", string(data)}, " ")
}
