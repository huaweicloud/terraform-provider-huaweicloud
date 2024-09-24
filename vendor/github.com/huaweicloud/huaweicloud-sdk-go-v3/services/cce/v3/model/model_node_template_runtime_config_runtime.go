package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NodeTemplateRuntimeConfigRuntime struct {
	Name *string `json:"name,omitempty"`
}

func (o NodeTemplateRuntimeConfigRuntime) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeTemplateRuntimeConfigRuntime struct{}"
	}

	return strings.Join([]string{"NodeTemplateRuntimeConfigRuntime", string(data)}, " ")
}
