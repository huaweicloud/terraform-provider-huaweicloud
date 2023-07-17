package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSnapshotConfigResponse Response Object
type DeleteSnapshotConfigResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteSnapshotConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSnapshotConfigResponse struct{}"
	}

	return strings.Join([]string{"DeleteSnapshotConfigResponse", string(data)}, " ")
}
