package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetInstancesProxyRestartResponse Response Object
type SetInstancesProxyRestartResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SetInstancesProxyRestartResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetInstancesProxyRestartResponse struct{}"
	}

	return strings.Join([]string{"SetInstancesProxyRestartResponse", string(data)}, " ")
}
