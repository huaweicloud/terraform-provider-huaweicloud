package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopSnapshotResponse Response Object
type StopSnapshotResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StopSnapshotResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopSnapshotResponse struct{}"
	}

	return strings.Join([]string{"StopSnapshotResponse", string(data)}, " ")
}
