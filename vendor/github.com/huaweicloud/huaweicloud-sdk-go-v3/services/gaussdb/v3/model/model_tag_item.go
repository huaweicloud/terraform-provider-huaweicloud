package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TagItem struct {
	// 标签键。最大长度36个unicode字符，不能为null或者空字符串，不能为空格。 字符集：0-9，A-Z，a-z，“_”，“-”，中文。

	Key string `json:"key"`
	// 标签值。最大长度43个unicode字符，可以为空字符串，不能为空格。 字符集：0-9，A-Z，a-z，“_”，“.”，“-”，中文。 - “action”值为“create”时，该参数必选。 - “action”值为“delete”时，如果value有值，按照key/value删除，如果value没值，则按照key删除。

	Value *string `json:"value,omitempty"`
}

func (o TagItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagItem struct{}"
	}

	return strings.Join([]string{"TagItem", string(data)}, " ")
}
