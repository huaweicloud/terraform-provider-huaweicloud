package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyInstanceConfigsResponse Response Object
type ModifyInstanceConfigsResponse struct {

	// 配置修改任务ID。
	JobId *string `json:"job_id,omitempty"`

	// 本次修改动态配置参数个数。
	DynamicConfig *int32 `json:"dynamic_config,omitempty"`

	// 本次修改静态配置参数个数。
	StaticConfig   *int32 `json:"static_config,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ModifyInstanceConfigsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyInstanceConfigsResponse struct{}"
	}

	return strings.Join([]string{"ModifyInstanceConfigsResponse", string(data)}, " ")
}
