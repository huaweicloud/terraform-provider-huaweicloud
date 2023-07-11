package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AutoLaunchStatisticsResponseInfo 自启动项统计信息
type AutoLaunchStatisticsResponseInfo struct {

	// 自启动项名称
	Name *string `json:"name,omitempty"`

	// 自启动项类型
	Type *string `json:"type,omitempty"`

	// 数量
	Num *int32 `json:"num,omitempty"`
}

func (o AutoLaunchStatisticsResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutoLaunchStatisticsResponseInfo struct{}"
	}

	return strings.Join([]string{"AutoLaunchStatisticsResponseInfo", string(data)}, " ")
}
