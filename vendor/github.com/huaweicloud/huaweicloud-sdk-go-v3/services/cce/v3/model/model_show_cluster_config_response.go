package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClusterConfigResponse Response Object
type ShowClusterConfigResponse struct {

	// 存储时长
	TtlInDays *int32 `json:"ttl_in_days,omitempty"`

	// 日志配置项
	LogConfigs     *[]ClusterLogConfigLogConfigs `json:"log_configs,omitempty"`
	HttpStatusCode int                           `json:"-"`
}

func (o ShowClusterConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterConfigResponse struct{}"
	}

	return strings.Join([]string{"ShowClusterConfigResponse", string(data)}, " ")
}
