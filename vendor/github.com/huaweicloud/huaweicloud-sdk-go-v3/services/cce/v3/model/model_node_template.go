package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NodeTemplate struct {
	Os *string `json:"os,omitempty"`

	ImageID *string `json:"imageID,omitempty"`

	Login *NodeTemplateLogin `json:"login,omitempty"`

	LifeCycle *NodeTemplateLifeCycle `json:"lifeCycle,omitempty"`

	RuntimeConfig *NodeTemplateRuntimeConfig `json:"runtimeConfig,omitempty"`

	ExtendParam *NodeTemplateExtendParam `json:"extendParam,omitempty"`
}

func (o NodeTemplate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeTemplate struct{}"
	}

	return strings.Join([]string{"NodeTemplate", string(data)}, " ")
}
