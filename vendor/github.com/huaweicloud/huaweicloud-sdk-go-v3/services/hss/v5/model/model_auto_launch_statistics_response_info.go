package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AutoLaunchStatisticsResponseInfo 自启动项统计信息
type AutoLaunchStatisticsResponseInfo struct {

	// 自启动项名称
	Name *string `json:"name,omitempty"`

	// 自启动项类型   - 0 ：自启动服务   - 1 ：定时任务   - 2 ：预加载动态库   - 3 ：Run注册表键   - 4 ：开机启动文件夹
	Type *string `json:"type,omitempty"`

	// 当前自启动项的主机数量
	Num *int32 `json:"num,omitempty"`
}

func (o AutoLaunchStatisticsResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutoLaunchStatisticsResponseInfo struct{}"
	}

	return strings.Join([]string{"AutoLaunchStatisticsResponseInfo", string(data)}, " ")
}
