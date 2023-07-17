package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HostNum 关联服务器数
type HostNum struct {
}

func (o HostNum) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HostNum struct{}"
	}

	return strings.Join([]string{"HostNum", string(data)}, " ")
}
