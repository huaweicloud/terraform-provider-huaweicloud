package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListClusterMasterSnapshotTasksResponse Response Object
type ListClusterMasterSnapshotTasksResponse struct {

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

func (o ListClusterMasterSnapshotTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClusterMasterSnapshotTasksResponse struct{}"
	}

	return strings.Join([]string{"ListClusterMasterSnapshotTasksResponse", string(data)}, " ")
}
