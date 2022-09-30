package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type StopAutoCreateSnapshotsRequest struct {

	// 快照所属的集群的ID。
	ClusterId string `json:"cluster_id"`
}

func (o StopAutoCreateSnapshotsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopAutoCreateSnapshotsRequest struct{}"
	}

	return strings.Join([]string{"StopAutoCreateSnapshotsRequest", string(data)}, " ")
}
