package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopSyncTaskResponse Response Object
type StopSyncTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StopSyncTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopSyncTaskResponse struct{}"
	}

	return strings.Join([]string{"StopSyncTaskResponse", string(data)}, " ")
}
