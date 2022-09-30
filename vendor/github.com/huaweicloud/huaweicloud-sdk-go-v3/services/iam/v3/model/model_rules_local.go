package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 联邦用户映射到IAM中的身份信息
type RulesLocal struct {
	User *RulesLocalUser `json:"user,omitempty"`

	Group *RulesLocalGroup `json:"group,omitempty"`

	Groups *RulesLocalGroups `json:"groups,omitempty"`
}

func (o RulesLocal) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RulesLocal struct{}"
	}

	return strings.Join([]string{"RulesLocal", string(data)}, " ")
}
