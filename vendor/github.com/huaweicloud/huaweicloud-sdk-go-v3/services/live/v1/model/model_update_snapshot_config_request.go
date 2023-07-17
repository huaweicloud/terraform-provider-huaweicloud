package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateSnapshotConfigRequest Request Object
type UpdateSnapshotConfigRequest struct {
	Body *LiveSnapshotConfig `json:"body,omitempty"`
}

func (o UpdateSnapshotConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSnapshotConfigRequest struct{}"
	}

	return strings.Join([]string{"UpdateSnapshotConfigRequest", string(data)}, " ")
}
