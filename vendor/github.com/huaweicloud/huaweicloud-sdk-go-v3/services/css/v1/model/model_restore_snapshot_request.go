package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RestoreSnapshotRequest Request Object
type RestoreSnapshotRequest struct {

	// 恢复快照所属的集群ID。
	ClusterId string `json:"cluster_id"`

	// 快照ID。
	SnapshotId string `json:"snapshot_id"`

	Body *RestoreSnapshotReq `json:"body,omitempty"`
}

func (o RestoreSnapshotRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreSnapshotRequest struct{}"
	}

	return strings.Join([]string{"RestoreSnapshotRequest", string(data)}, " ")
}
