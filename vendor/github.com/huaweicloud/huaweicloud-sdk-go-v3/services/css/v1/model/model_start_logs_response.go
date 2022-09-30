package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type StartLogsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StartLogsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartLogsResponse struct{}"
	}

	return strings.Join([]string{"StartLogsResponse", string(data)}, " ")
}
