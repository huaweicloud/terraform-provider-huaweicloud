package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResourceTagItem struct {
	// 标签键。

	Key string `json:"key"`
	// 标签值。

	Value string `json:"value"`
}

func (o ResourceTagItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceTagItem struct{}"
	}

	return strings.Join([]string{"ResourceTagItem", string(data)}, " ")
}
