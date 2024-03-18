package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SetHostPrivilegeRequestV3 struct {

	// 数据库用户名
	UserName string `json:"user_name"`

	// host信息
	Hosts []string `json:"hosts"`
}

func (o SetHostPrivilegeRequestV3) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetHostPrivilegeRequestV3 struct{}"
	}

	return strings.Join([]string{"SetHostPrivilegeRequestV3", string(data)}, " ")
}
