package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type TokenRole struct {

	// 权限名称。
	Name string `json:"name"`

	// 权限ID。默认显示为0，非真实权限ID。
	Id string `json:"id"`
}

func (o TokenRole) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TokenRole struct{}"
	}

	return strings.Join([]string{"TokenRole", string(data)}, " ")
}
