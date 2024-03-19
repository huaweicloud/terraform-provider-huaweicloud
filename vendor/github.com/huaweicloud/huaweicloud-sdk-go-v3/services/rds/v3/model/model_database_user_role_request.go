package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DatabaseUserRoleRequest struct {

	// 用户名称
	User string `json:"user"`

	// 角色名称
	Roles []string `json:"roles"`
}

func (o DatabaseUserRoleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DatabaseUserRoleRequest struct{}"
	}

	return strings.Join([]string{"DatabaseUserRoleRequest", string(data)}, " ")
}
