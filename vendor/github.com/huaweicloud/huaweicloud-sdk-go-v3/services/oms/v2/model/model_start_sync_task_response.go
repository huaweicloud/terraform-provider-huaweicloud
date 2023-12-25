package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartSyncTaskResponse Response Object
type StartSyncTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StartSyncTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartSyncTaskResponse struct{}"
	}

	return strings.Join([]string{"StartSyncTaskResponse", string(data)}, " ")
}
