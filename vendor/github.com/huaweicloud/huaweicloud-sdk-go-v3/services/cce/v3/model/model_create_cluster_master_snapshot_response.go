package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClusterMasterSnapshotResponse Response Object
type CreateClusterMasterSnapshotResponse struct {

	// 任务ID
	Uid *string `json:"uid,omitempty"`

	Metadata       *SnapshotCluserResponseMetadata `json:"metadata,omitempty"`
	HttpStatusCode int                             `json:"-"`
}

func (o CreateClusterMasterSnapshotResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterMasterSnapshotResponse struct{}"
	}

	return strings.Join([]string{"CreateClusterMasterSnapshotResponse", string(data)}, " ")
}
