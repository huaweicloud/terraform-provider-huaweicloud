package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 登录源IP
type LoginIp struct {
}

func (o LoginIp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LoginIp struct{}"
	}

	return strings.Join([]string{"LoginIp", string(data)}, " ")
}
