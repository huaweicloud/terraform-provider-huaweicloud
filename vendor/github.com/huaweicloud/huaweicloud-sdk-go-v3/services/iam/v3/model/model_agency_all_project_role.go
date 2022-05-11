package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type AgencyAllProjectRole struct {

	// 权限ID。
	Id string `json:"id"`

	Links *LinksSelf `json:"links"`

	// 权限名。
	Name string `json:"name"`
}

func (o AgencyAllProjectRole) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AgencyAllProjectRole struct{}"
	}

	return strings.Join([]string{"AgencyAllProjectRole", string(data)}, " ")
}
