package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 联邦用户在本系统中的用户名称
type RulesLocalUser struct {

	// 联邦用户在本系统中的用户名称
	Name string `json:"name"`
}

func (o RulesLocalUser) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RulesLocalUser struct{}"
	}

	return strings.Join([]string{"RulesLocalUser", string(data)}, " ")
}
