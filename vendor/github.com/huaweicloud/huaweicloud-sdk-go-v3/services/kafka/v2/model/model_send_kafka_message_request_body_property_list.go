package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SendKafkaMessageRequestBodyPropertyList struct {

	// 特性名字
	Name *string `json:"name,omitempty"`

	// 特性值
	Value *string `json:"value,omitempty"`
}

func (o SendKafkaMessageRequestBodyPropertyList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SendKafkaMessageRequestBodyPropertyList struct{}"
	}

	return strings.Join([]string{"SendKafkaMessageRequestBodyPropertyList", string(data)}, " ")
}
