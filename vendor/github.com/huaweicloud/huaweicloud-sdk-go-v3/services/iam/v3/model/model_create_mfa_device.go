package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// {  \"virtual_mfa_device\": {   \"name\": \"{divice_name}\",   \"user_id\": \"{user_id}\"  } }
type CreateMfaDevice struct {

	// 设备名称。 最小长度：1 最大长度：64
	Name string `json:"name"`

	// 创建MFA设备的IAM用户ID。
	UserId string `json:"user_id"`
}

func (o CreateMfaDevice) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMfaDevice struct{}"
	}

	return strings.Join([]string{"CreateMfaDevice", string(data)}, " ")
}
