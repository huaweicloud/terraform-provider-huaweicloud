package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type StartAutoCreateSnapshotsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StartAutoCreateSnapshotsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartAutoCreateSnapshotsResponse struct{}"
	}

	return strings.Join([]string{"StartAutoCreateSnapshotsResponse", string(data)}, " ")
}
