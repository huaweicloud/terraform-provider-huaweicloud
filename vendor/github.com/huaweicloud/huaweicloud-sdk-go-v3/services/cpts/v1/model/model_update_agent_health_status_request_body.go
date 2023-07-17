package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAgentHealthStatusRequestBody 上报探针健康状态请求体
type UpdateAgentHealthStatusRequestBody struct {

	// 探针上次获取全链路应用的更新时间戳（单位：毫秒）
	UpdateTime int64 `json:"update_time"`
}

func (o UpdateAgentHealthStatusRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAgentHealthStatusRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateAgentHealthStatusRequestBody", string(data)}, " ")
}
