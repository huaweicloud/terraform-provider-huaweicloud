package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ClusterLogConfig struct {

	// 存储时长
	TtlInDays *int32 `json:"ttl_in_days,omitempty"`

	// 日志配置项
	LogConfigs *[]ClusterLogConfigLogConfigs `json:"log_configs,omitempty"`
}

func (o ClusterLogConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterLogConfig struct{}"
	}

	return strings.Join([]string{"ClusterLogConfig", string(data)}, " ")
}
