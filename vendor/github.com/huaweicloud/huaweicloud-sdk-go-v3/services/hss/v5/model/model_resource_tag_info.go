package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResourceTagInfo struct {

	// 键。最大长度128个unicode字符。 key不能为空
	Key string `json:"key"`

	// 值
	Value string `json:"value"`
}

func (o ResourceTagInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceTagInfo struct{}"
	}

	return strings.Join([]string{"ResourceTagInfo", string(data)}, " ")
}
