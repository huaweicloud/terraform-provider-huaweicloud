package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// LoginUserName 登录用户名
type LoginUserName struct {
}

func (o LoginUserName) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LoginUserName struct{}"
	}

	return strings.Join([]string{"LoginUserName", string(data)}, " ")
}
