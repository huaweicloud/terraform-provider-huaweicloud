package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NodeTemplateExtendParam struct {
	UserID *string `json:"userID,omitempty"`
}

func (o NodeTemplateExtendParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeTemplateExtendParam struct{}"
	}

	return strings.Join([]string{"NodeTemplateExtendParam", string(data)}, " ")
}
