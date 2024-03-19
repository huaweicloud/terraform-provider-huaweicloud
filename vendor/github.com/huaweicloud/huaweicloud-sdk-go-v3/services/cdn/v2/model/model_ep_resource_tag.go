package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EpResourceTag 标签信息。
type EpResourceTag struct {

	// 资源标签key。
	Key string `json:"key"`

	// 资源标签value值。
	Value string `json:"value"`
}

func (o EpResourceTag) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EpResourceTag struct{}"
	}

	return strings.Join([]string{"EpResourceTag", string(data)}, " ")
}
