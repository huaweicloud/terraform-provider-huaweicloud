package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneAddUserToGroupRequest struct {

	// 用户组ID，获取方式请参见：[获取用户组ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	GroupId string `json:"group_id"`

	// 待添加的IAM用户ID，获取方式请参见：[获取IAM用户ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	UserId string `json:"user_id"`
}

func (o KeystoneAddUserToGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneAddUserToGroupRequest struct{}"
	}

	return strings.Join([]string{"KeystoneAddUserToGroupRequest", string(data)}, " ")
}
