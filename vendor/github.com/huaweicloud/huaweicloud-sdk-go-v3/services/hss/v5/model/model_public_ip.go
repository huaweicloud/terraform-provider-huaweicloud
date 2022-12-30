package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 弹性公网IP地址
type PublicIp struct {
}

func (o PublicIp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PublicIp struct{}"
	}

	return strings.Join([]string{"PublicIp", string(data)}, " ")
}
