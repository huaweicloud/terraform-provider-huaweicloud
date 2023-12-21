package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResourceDeleteTag CCE资源标签，用于唯一标识待删除的资源标签
type ResourceDeleteTag struct {

	// Key值。 - 不能为空，最多支持128个字符 - 可用UTF-8格式表示的汉字、字母、数字和空格 - 支持部分特殊字符：_.:/=+-@ - 不能以\"\\_sys\\_\"开头
	Key *string `json:"key,omitempty"`
}

func (o ResourceDeleteTag) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceDeleteTag struct{}"
	}

	return strings.Join([]string{"ResourceDeleteTag", string(data)}, " ")
}
