package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProjectTagItem struct {
	// 标签键。

	Key string `json:"key"`
	// 标签值。

	Values []string `json:"values"`
}

func (o ProjectTagItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProjectTagItem struct{}"
	}

	return strings.Join([]string{"ProjectTagItem", string(data)}, " ")
}
