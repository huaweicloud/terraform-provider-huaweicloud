package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ClusterLogConfigLogConfigs struct {

	// 日志类型
	Name *string `json:"name,omitempty"`

	// 是否采集
	Enable *bool `json:"enable,omitempty"`
}

func (o ClusterLogConfigLogConfigs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterLogConfigLogConfigs struct{}"
	}

	return strings.Join([]string{"ClusterLogConfigLogConfigs", string(data)}, " ")
}
