package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateSyncEventsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateSyncEventsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSyncEventsResponse struct{}"
	}

	return strings.Join([]string{"CreateSyncEventsResponse", string(data)}, " ")
}
