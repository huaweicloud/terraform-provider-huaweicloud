package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSyncTaskResponse Response Object
type DeleteSyncTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteSyncTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSyncTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteSyncTaskResponse", string(data)}, " ")
}
