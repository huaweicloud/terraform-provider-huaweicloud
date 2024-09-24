package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateAutopilotClusterMasterSnapshotResponse Response Object
type CreateAutopilotClusterMasterSnapshotResponse struct {

	// 任务ID
	Uid *string `json:"uid,omitempty"`

	Metadata       *SnapshotCluserResponseMetadata `json:"metadata,omitempty"`
	HttpStatusCode int                             `json:"-"`
}

func (o CreateAutopilotClusterMasterSnapshotResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAutopilotClusterMasterSnapshotResponse struct{}"
	}

	return strings.Join([]string{"CreateAutopilotClusterMasterSnapshotResponse", string(data)}, " ")
}
