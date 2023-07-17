package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSnapshotRequest Request Object
type DeleteSnapshotRequest struct {

	// 删除快照所属的集群的ID。
	ClusterId string `json:"cluster_id"`

	// 要删除快照的ID。
	SnapshotId string `json:"snapshot_id"`
}

func (o DeleteSnapshotRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSnapshotRequest struct{}"
	}

	return strings.Join([]string{"DeleteSnapshotRequest", string(data)}, " ")
}
