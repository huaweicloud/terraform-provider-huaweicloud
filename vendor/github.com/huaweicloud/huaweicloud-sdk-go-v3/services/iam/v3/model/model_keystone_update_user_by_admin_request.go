package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneUpdateUserByAdminRequest struct {

	// 待修改信息的IAM用户ID，获取方式请参见：[获取用户ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	UserId string `json:"user_id"`

	Body *KeystoneUpdateUserByAdminRequestBody `json:"body,omitempty"`
}

func (o KeystoneUpdateUserByAdminRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneUpdateUserByAdminRequest struct{}"
	}

	return strings.Join([]string{"KeystoneUpdateUserByAdminRequest", string(data)}, " ")
}
