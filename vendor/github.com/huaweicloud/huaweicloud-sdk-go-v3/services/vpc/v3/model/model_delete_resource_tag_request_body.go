package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeleteResourceTagRequestBody struct {

	// 功能说明：标签键 约束：同一资源的key值不能重复。
	Key string `json:"key"`

	// 功能说明：标签值
	Value *string `json:"value,omitempty"`
}

func (o DeleteResourceTagRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteResourceTagRequestBody struct{}"
	}

	return strings.Join([]string{"DeleteResourceTagRequestBody", string(data)}, " ")
}
