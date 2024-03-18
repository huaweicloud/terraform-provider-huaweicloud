package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPreCheckTasksResponse Response Object
type ListPreCheckTasksResponse struct {

	// api版本，默认为v3
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 类型
	Kind *string `json:"kind,omitempty"`

	Metadata *Metadata `json:"metadata,omitempty"`

	// 集群检查任务列表
	Items          *[]PrecheckClusterTask `json:"items,omitempty"`
	HttpStatusCode int                    `json:"-"`
}

func (o ListPreCheckTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPreCheckTasksResponse struct{}"
	}

	return strings.Join([]string{"ListPreCheckTasksResponse", string(data)}, " ")
}
