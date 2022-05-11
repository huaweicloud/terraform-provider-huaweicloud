package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type Tag struct {

	// 1.功能说明：标签键 2.取值范围：最大长度36个unicode字符。 key不能为空。不能包含非打印字符ASCII(0-31)，*,<,>,\\,=
	Key string `json:"key"`

	// 1. 功能描述：标签值 2. 取值范围：每个值最大长度43个unicode字符，可以为空字符串。 不能包含非打印字符ASCII(0-31)，*,<,>,\\,=
	Value string `json:"value"`
}

func (o Tag) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Tag struct{}"
	}

	return strings.Join([]string{"Tag", string(data)}, " ")
}
