package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateSnapshotConfigResponse Response Object
type CreateSnapshotConfigResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateSnapshotConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSnapshotConfigResponse struct{}"
	}

	return strings.Join([]string{"CreateSnapshotConfigResponse", string(data)}, " ")
}
