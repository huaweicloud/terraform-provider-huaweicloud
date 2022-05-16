package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateCredentialResult struct {

	// IAM用户ID。
	UserId string `json:"user_id"`

	// 修改的AK。
	Access string `json:"access"`

	// 访问密钥状态。
	Status string `json:"status"`

	// 访问密钥创建时间。
	CreateTime string `json:"create_time"`

	// 访问密钥描述信息。
	Description string `json:"description"`
}

func (o UpdateCredentialResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCredentialResult struct{}"
	}

	return strings.Join([]string{"UpdateCredentialResult", string(data)}, " ")
}
