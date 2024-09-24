package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AutopilotResourceTag CCE资源标签
type AutopilotResourceTag struct {

	// Key值。 - 不能为空，最多支持128个字符 - 可用UTF-8格式表示的汉字、字母、数字和空格 - 支持部分特殊字符：_.:/=+-@ - 不能以\"\\_sys\\_\"开头
	Key *string `json:"key,omitempty"`

	// Value值。 - 可以为空但不能缺省，最多支持255个字符 - 可用UTF-8格式表示的汉字、字母、数字和空格 - 支持部分特殊字符：_.:/=+-@
	Value *string `json:"value,omitempty"`
}

func (o AutopilotResourceTag) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotResourceTag struct{}"
	}

	return strings.Join([]string{"AutopilotResourceTag", string(data)}, " ")
}
