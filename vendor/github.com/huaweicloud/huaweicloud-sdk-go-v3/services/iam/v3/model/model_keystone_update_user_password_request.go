package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneUpdateUserPasswordRequest struct {

	// 待修改密码的IAM用户ID，获取方式请参见：[获取用户ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	UserId string `json:"user_id"`

	Body *KeystoneUpdateUserPasswordRequestBody `json:"body,omitempty"`
}

func (o KeystoneUpdateUserPasswordRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneUpdateUserPasswordRequest struct{}"
	}

	return strings.Join([]string{"KeystoneUpdateUserPasswordRequest", string(data)}, " ")
}
