package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NodeTemplateLifeCycle struct {
	PreInstall *string `json:"preInstall,omitempty"`

	PostInstall *string `json:"postInstall,omitempty"`
}

func (o NodeTemplateLifeCycle) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeTemplateLifeCycle struct{}"
	}

	return strings.Join([]string{"NodeTemplateLifeCycle", string(data)}, " ")
}
