package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 通知用户列表。
type NotificationUsers struct {

	// IAM用户组。
	UserGroup string `json:"user_group"`

	// IAM用户。
	UserList []string `json:"user_list"`
}

func (o NotificationUsers) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NotificationUsers struct{}"
	}

	return strings.Join([]string{"NotificationUsers", string(data)}, " ")
}
