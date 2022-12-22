package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 进程统计信息
type ProcessStatisticResponseInfo struct {

	// 进程名称
	Path *string `json:"path,omitempty"`

	// 进程数量
	Num *int32 `json:"num,omitempty"`
}

func (o ProcessStatisticResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProcessStatisticResponseInfo struct{}"
	}

	return strings.Join([]string{"ProcessStatisticResponseInfo", string(data)}, " ")
}
