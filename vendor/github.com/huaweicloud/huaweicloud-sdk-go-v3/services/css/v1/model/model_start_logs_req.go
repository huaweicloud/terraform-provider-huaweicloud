package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StartLogsReq struct {

	// IAM委托。
	Agency string `json:"agency"`

	// 日志在OBS桶中的备份路径。
	LogBasePath string `json:"logBasePath"`

	// 用于存储日志的OBS桶的桶名。
	LogBucket string `json:"logBucket"`
}

func (o StartLogsReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartLogsReq struct{}"
	}

	return strings.Join([]string{"StartLogsReq", string(data)}, " ")
}
