package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NodeTemplateLoginUserPassword struct {
	Username *string `json:"username,omitempty"`

	Password *string `json:"password,omitempty"`
}

func (o NodeTemplateLoginUserPassword) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeTemplateLoginUserPassword struct{}"
	}

	return strings.Join([]string{"NodeTemplateLoginUserPassword", string(data)}, " ")
}
