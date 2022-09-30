package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateLogSettingRequest struct {

	// 指定更改日志基础配置的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *UpdateLogSettingReq `json:"body,omitempty"`
}

func (o UpdateLogSettingRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateLogSettingRequest struct{}"
	}

	return strings.Join([]string{"UpdateLogSettingRequest", string(data)}, " ")
}
