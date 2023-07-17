package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopLogsResponse Response Object
type StopLogsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StopLogsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopLogsResponse struct{}"
	}

	return strings.Join([]string{"StopLogsResponse", string(data)}, " ")
}
