package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 进程统计信息
type AppStatisticResponseInfo struct {

	// 软件名称
	AppName *string `json:"app_name,omitempty"`

	// 进程数量
	Num *int32 `json:"num,omitempty"`
}

func (o AppStatisticResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AppStatisticResponseInfo struct{}"
	}

	return strings.Join([]string{"AppStatisticResponseInfo", string(data)}, " ")
}
