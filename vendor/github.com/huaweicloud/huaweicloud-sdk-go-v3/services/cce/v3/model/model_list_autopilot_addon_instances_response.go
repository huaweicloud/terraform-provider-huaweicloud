package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotAddonInstancesResponse Response Object
type ListAutopilotAddonInstancesResponse struct {

	// API类型，固定值“Addon”，该值不可修改。
	Kind *string `json:"kind,omitempty"`

	// API版本，固定值“v3”，该值不可修改。
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 插件实例列表
	Items          *[]AddonInstance `json:"items,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ListAutopilotAddonInstancesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotAddonInstancesResponse struct{}"
	}

	return strings.Join([]string{"ListAutopilotAddonInstancesResponse", string(data)}, " ")
}
