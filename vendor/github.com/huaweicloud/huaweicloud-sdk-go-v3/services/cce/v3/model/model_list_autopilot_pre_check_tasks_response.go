package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotPreCheckTasksResponse Response Object
type ListAutopilotPreCheckTasksResponse struct {

	// api版本，默认为v3
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 类型
	Kind *string `json:"kind,omitempty"`

	Metadata *Metadata `json:"metadata,omitempty"`

	// 集群检查任务列表
	Items          *[]PrecheckClusterTask `json:"items,omitempty"`
	HttpStatusCode int                    `json:"-"`
}

func (o ListAutopilotPreCheckTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotPreCheckTasksResponse struct{}"
	}

	return strings.Join([]string{"ListAutopilotPreCheckTasksResponse", string(data)}, " ")
}
