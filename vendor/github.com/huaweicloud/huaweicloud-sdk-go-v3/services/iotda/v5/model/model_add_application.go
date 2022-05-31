package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 创建资源空间结构体。
type AddApplication struct {

	// **参数说明**：资源空间名称。 **取值范围**：长度不超过64，允许中文、字母、数字、下划线（_）、连接符（-）等一些特殊字符的组合。
	AppName string `json:"app_name"`
}

func (o AddApplication) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddApplication struct{}"
	}

	return strings.Join([]string{"AddApplication", string(data)}, " ")
}
