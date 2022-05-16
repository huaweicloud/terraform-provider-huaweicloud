package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type CreateCredentialResult struct {

	// 创建访问密钥时间。
	CreateTime string `json:"create_time"`

	// 创建的AK。
	Access string `json:"access"`

	// 创建的SK。
	Secret string `json:"secret"`

	// 访问密钥状态。
	Status string `json:"status"`

	// IAM用户ID。
	UserId string `json:"user_id"`

	// 访问密钥描述信息。
	Description string `json:"description"`
}

func (o CreateCredentialResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCredentialResult struct{}"
	}

	return strings.Join([]string{"CreateCredentialResult", string(data)}, " ")
}
