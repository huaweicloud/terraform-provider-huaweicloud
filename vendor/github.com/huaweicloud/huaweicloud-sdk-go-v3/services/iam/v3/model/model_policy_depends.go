package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type PolicyDepends struct {

	// 权限所在目录。
	Catalog string `json:"catalog"`

	// 权限展示名。
	DisplayName string `json:"display_name"`
}

func (o PolicyDepends) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PolicyDepends struct{}"
	}

	return strings.Join([]string{"PolicyDepends", string(data)}, " ")
}
