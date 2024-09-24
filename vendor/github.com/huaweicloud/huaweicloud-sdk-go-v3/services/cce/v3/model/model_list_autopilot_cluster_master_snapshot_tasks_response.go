package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotClusterMasterSnapshotTasksResponse Response Object
type ListAutopilotClusterMasterSnapshotTasksResponse struct {

	// api版本，默认为v3.1
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 任务类型
	Kind *string `json:"kind,omitempty"`

	Metadata *SnapshotTaskMetadata `json:"metadata,omitempty"`

	// 备份任务列表
	Items *[]SnapshotTask `json:"items,omitempty"`

	Status         *SnapshotTaskStatus `json:"status,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ListAutopilotClusterMasterSnapshotTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotClusterMasterSnapshotTasksResponse struct{}"
	}

	return strings.Join([]string{"ListAutopilotClusterMasterSnapshotTasksResponse", string(data)}, " ")
}
