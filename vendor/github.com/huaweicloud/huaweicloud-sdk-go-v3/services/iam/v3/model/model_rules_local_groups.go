package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 联邦用户在本系统中所属用户组列表
type RulesLocalGroups struct {

	// 联邦用户在本系统中所属用户组列表
	Name string `json:"name"`
}

func (o RulesLocalGroups) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RulesLocalGroups struct{}"
	}

	return strings.Join([]string{"RulesLocalGroups", string(data)}, " ")
}
