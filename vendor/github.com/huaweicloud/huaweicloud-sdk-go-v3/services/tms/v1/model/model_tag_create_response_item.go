package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TagCreateResponseItem struct {

	// 资源ID。
	ResourceId string `json:"resource_id"`

	// 资源类型。
	ResourceType string `json:"resource_type"`

	// 错误码
	ErrorCode string `json:"error_code"`

	// 错误描述
	ErrorMsg string `json:"error_msg"`
}

func (o TagCreateResponseItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagCreateResponseItem struct{}"
	}

	return strings.Join([]string{"TagCreateResponseItem", string(data)}, " ")
}
