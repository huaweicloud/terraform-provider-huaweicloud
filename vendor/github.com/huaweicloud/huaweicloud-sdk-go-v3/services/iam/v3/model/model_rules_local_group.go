package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 联邦用户在本系统中所属用户组
type RulesLocalGroup struct {

	// 联邦用户在本系统中所属用户组
	Name string `json:"name"`
}

func (o RulesLocalGroup) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RulesLocalGroup struct{}"
	}

	return strings.Join([]string{"RulesLocalGroup", string(data)}, " ")
}
