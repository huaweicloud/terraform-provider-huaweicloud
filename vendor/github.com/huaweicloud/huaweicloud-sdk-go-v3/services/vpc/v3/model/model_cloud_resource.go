package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type CloudResource struct {

	// 功能描述：资源类型
	ResourceType string `json:"resource_type"`

	// 功能说明：资源数量
	ResourceCount int32 `json:"resource_count"`
}

func (o CloudResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CloudResource struct{}"
	}

	return strings.Join([]string{"CloudResource", string(data)}, " ")
}
