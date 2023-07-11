package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateSnapshotConfigRequest Request Object
type CreateSnapshotConfigRequest struct {
	Body *LiveSnapshotConfig `json:"body,omitempty"`
}

func (o CreateSnapshotConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSnapshotConfigRequest struct{}"
	}

	return strings.Join([]string{"CreateSnapshotConfigRequest", string(data)}, " ")
}
