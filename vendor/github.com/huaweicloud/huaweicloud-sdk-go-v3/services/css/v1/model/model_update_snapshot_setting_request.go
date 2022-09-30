package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateSnapshotSettingRequest struct {

	// 指定待操作的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *UpdateSnapshotSettingReq `json:"body,omitempty"`
}

func (o UpdateSnapshotSettingRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSnapshotSettingRequest struct{}"
	}

	return strings.Join([]string{"UpdateSnapshotSettingRequest", string(data)}, " ")
}
