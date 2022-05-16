package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateLoginProtectRequest struct {

	// 待修改登录保护状态信息的IAM用户ID。
	UserId string `json:"user_id"`

	Body *UpdateLoginProjectReq `json:"body,omitempty"`
}

func (o UpdateLoginProtectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateLoginProtectRequest struct{}"
	}

	return strings.Join([]string{"UpdateLoginProtectRequest", string(data)}, " ")
}
