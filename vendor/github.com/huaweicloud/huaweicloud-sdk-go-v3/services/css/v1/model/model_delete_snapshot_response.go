package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSnapshotResponse Response Object
type DeleteSnapshotResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteSnapshotResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSnapshotResponse struct{}"
	}

	return strings.Join([]string{"DeleteSnapshotResponse", string(data)}, " ")
}
