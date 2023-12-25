package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AddLogConfigResponseBody struct {

	// 实例日志配置信息。
	LogConfigs []AddLogConfigs `json:"log_configs"`
}

func (o AddLogConfigResponseBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddLogConfigResponseBody struct{}"
	}

	return strings.Join([]string{"AddLogConfigResponseBody", string(data)}, " ")
}
