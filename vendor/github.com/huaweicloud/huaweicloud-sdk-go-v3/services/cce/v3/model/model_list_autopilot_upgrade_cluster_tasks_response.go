package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotUpgradeClusterTasksResponse Response Object
type ListAutopilotUpgradeClusterTasksResponse struct {

	// api版本，默认为v3
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 资源类型
	Kind *string `json:"kind,omitempty"`

	Metadata *UpgradeTaskMetadata `json:"metadata,omitempty"`

	// 集群升级任务列表
	Items          *[]UpgradeTaskResponseBody `json:"items,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o ListAutopilotUpgradeClusterTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotUpgradeClusterTasksResponse struct{}"
	}

	return strings.Join([]string{"ListAutopilotUpgradeClusterTasksResponse", string(data)}, " ")
}
