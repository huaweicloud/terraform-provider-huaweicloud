package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type StopSnapshotRequest struct {

	// 停用快照所属的集群的ID。
	ClusterId string `json:"cluster_id"`
}

func (o StopSnapshotRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopSnapshotRequest struct{}"
	}

	return strings.Join([]string{"StopSnapshotRequest", string(data)}, " ")
}
