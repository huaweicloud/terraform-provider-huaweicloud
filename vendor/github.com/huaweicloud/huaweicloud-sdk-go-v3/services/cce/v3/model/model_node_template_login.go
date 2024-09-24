package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NodeTemplateLogin struct {
	SshKey *string `json:"sshKey,omitempty"`

	UserPassword *NodeTemplateLoginUserPassword `json:"userPassword,omitempty"`
}

func (o NodeTemplateLogin) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeTemplateLogin struct{}"
	}

	return strings.Join([]string{"NodeTemplateLogin", string(data)}, " ")
}
