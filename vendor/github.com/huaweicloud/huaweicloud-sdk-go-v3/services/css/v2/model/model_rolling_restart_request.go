package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RollingRestartRequest Request Object
type RollingRestartRequest struct {

	// 指定重启集群ID。
	ClusterId string `json:"cluster_id"`

	Body *RollingRestartReq `json:"body,omitempty"`
}

func (o RollingRestartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RollingRestartRequest struct{}"
	}

	return strings.Join([]string{"RollingRestartRequest", string(data)}, " ")
}
