package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type InstanceTagItem struct {
	// 标签键。

	Key string `json:"key"`
	// 标签值。

	Value string `json:"value"`
}

func (o InstanceTagItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InstanceTagItem struct{}"
	}

	return strings.Join([]string{"InstanceTagItem", string(data)}, " ")
}
