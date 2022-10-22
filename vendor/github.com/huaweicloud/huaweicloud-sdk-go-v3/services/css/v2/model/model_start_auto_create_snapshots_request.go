package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type StartAutoCreateSnapshotsRequest struct {

	// 快照所属的集群的ID。
	ClusterId string `json:"cluster_id"`

	Body *StartAutoCreateSnapshotsReq `json:"body,omitempty"`
}

func (o StartAutoCreateSnapshotsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartAutoCreateSnapshotsRequest struct{}"
	}

	return strings.Join([]string{"StartAutoCreateSnapshotsRequest", string(data)}, " ")
}
