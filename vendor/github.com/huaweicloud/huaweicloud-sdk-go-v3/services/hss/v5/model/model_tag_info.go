package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TagInfo struct {

	// 键。最大长度128个unicode字符。 key不能为空
	Key *string `json:"key,omitempty"`

	// 值。最大长度255个unicode字符。
	Value *string `json:"value,omitempty"`
}

func (o TagInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagInfo struct{}"
	}

	return strings.Join([]string{"TagInfo", string(data)}, " ")
}
