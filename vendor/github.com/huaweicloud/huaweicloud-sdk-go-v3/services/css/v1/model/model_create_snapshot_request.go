package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateSnapshotRequest Request Object
type CreateSnapshotRequest struct {

	// 指定要创建快照的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *CreateSnapshotReq `json:"body,omitempty"`
}

func (o CreateSnapshotRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSnapshotRequest struct{}"
	}

	return strings.Join([]string{"CreateSnapshotRequest", string(data)}, " ")
}
