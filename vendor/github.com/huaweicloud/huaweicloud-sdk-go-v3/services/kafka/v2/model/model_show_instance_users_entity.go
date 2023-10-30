package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowInstanceUsersEntity struct {

	// 用户名称。  由英文字符开头，只能由英文字母、数字、中划线、下划线组成，长度为4~64的字符。
	UserName *string `json:"user_name,omitempty"`

	// 用户描述。
	UserDesc *string `json:"user_desc,omitempty"`

	// 用户角色。
	Role *string `json:"role,omitempty"`

	// 是否为默认应用。
	DefaultApp *bool `json:"default_app,omitempty"`

	// 创建时间。
	CreatedTime *int64 `json:"created_time,omitempty"`
}

func (o ShowInstanceUsersEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceUsersEntity struct{}"
	}

	return strings.Join([]string{"ShowInstanceUsersEntity", string(data)}, " ")
}
