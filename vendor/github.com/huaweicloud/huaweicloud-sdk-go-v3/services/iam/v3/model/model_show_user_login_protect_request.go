package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowUserLoginProtectRequest struct {

	// 待查询的IAM用户ID，获取方式请参见：[获取用户ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	UserId string `json:"user_id"`
}

func (o ShowUserLoginProtectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowUserLoginProtectRequest struct{}"
	}

	return strings.Join([]string{"ShowUserLoginProtectRequest", string(data)}, " ")
}
