package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResourceTagBody struct {

	// 资源ID
	ResourceId string `json:"resource_id"`

	// 资源类型
	ResourceType string `json:"resource_type"`
}

func (o ResourceTagBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceTagBody struct{}"
	}

	return strings.Join([]string{"ResourceTagBody", string(data)}, " ")
}
