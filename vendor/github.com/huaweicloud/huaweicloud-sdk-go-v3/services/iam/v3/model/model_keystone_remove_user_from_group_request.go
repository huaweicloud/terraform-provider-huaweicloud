package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneRemoveUserFromGroupRequest struct {

	// 用户组ID，获取方式请参见：[获取用户组ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	GroupId string `json:"group_id"`

	// 待从用户组中移除的IAM用户ID，获取方式请参见：[获取IAM用户ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	UserId string `json:"user_id"`
}

func (o KeystoneRemoveUserFromGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneRemoveUserFromGroupRequest struct{}"
	}

	return strings.Join([]string{"KeystoneRemoveUserFromGroupRequest", string(data)}, " ")
}
