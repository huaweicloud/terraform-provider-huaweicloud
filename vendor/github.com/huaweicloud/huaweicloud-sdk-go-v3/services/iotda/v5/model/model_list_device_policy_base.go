package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDevicePolicyBase 策略信息结构体，返回基本信息。
type ListDevicePolicyBase struct {

	// **参数说明**：资源空间ID。
	AppId *string `json:"app_id,omitempty"`

	// 策略ID。
	PolicyId *string `json:"policy_id,omitempty"`

	// 策略名称。
	PolicyName *string `json:"policy_name,omitempty"`

	// 在物联网平台创建策略的时间。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	CreateTime *string `json:"create_time,omitempty"`

	// 在物联网平台更新策略的时间。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	UpdateTime *string `json:"update_time,omitempty"`
}

func (o ListDevicePolicyBase) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDevicePolicyBase struct{}"
	}

	return strings.Join([]string{"ListDevicePolicyBase", string(data)}, " ")
}
