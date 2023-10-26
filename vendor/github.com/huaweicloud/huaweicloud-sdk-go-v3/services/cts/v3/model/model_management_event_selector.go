package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ManagementEventSelector 管理类事件选择器。
type ManagementEventSelector struct {

	// 标识不转储的云服务名称。 目前只支持设置为KMS，表示屏蔽KMS服务的createDatakey事件。
	ExcludeService *[]string `json:"exclude_service,omitempty"`
}

func (o ManagementEventSelector) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ManagementEventSelector struct{}"
	}

	return strings.Join([]string{"ManagementEventSelector", string(data)}, " ")
}
