package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateSnapshotConfigResponse Response Object
type UpdateSnapshotConfigResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateSnapshotConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSnapshotConfigResponse struct{}"
	}

	return strings.Join([]string{"UpdateSnapshotConfigResponse", string(data)}, " ")
}
