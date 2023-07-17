package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RestoreSnapshotResponse Response Object
type RestoreSnapshotResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o RestoreSnapshotResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreSnapshotResponse struct{}"
	}

	return strings.Join([]string{"RestoreSnapshotResponse", string(data)}, " ")
}
