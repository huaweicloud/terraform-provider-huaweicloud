package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartInstanceReduceVolumeActionResponse Response Object
type StartInstanceReduceVolumeActionResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o StartInstanceReduceVolumeActionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartInstanceReduceVolumeActionResponse struct{}"
	}

	return strings.Join([]string{"StartInstanceReduceVolumeActionResponse", string(data)}, " ")
}
