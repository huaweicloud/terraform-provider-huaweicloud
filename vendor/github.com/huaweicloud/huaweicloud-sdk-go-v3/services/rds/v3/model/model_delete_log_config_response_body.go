package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeleteLogConfigResponseBody struct {

	// 实例日志配置信息。
	LogConfigs []DeleteLogConfigs `json:"log_configs"`
}

func (o DeleteLogConfigResponseBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteLogConfigResponseBody struct{}"
	}

	return strings.Join([]string{"DeleteLogConfigResponseBody", string(data)}, " ")
}
