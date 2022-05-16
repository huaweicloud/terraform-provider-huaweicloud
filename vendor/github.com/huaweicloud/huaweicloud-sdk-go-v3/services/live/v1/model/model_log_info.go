package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type LogInfo struct {

	// 日志文件名，打包文件名格式：{Domain}_{logStartTimeStamp}.log.gz
	Name string `json:"name"`

	// 日志下载链接
	Url string `json:"url"`

	// 日志文件大小
	Size int64 `json:"size"`

	// 日志文件中日志开始时间，北京时间
	StartTime string `json:"start_time"`

	// 日志文件中日志结束时间，北京时间
	EndTime string `json:"end_time"`
}

func (o LogInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LogInfo struct{}"
	}

	return strings.Join([]string{"LogInfo", string(data)}, " ")
}
