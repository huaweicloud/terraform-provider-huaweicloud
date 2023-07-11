package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResetPasswordReq struct {

	// 安全模式下集群管理员admin的密码，只有在创建集群时authorityEnable设置为true时该接口可用。  - 参数范围：8~32个字符。  - 参数要求：密码至少包含大写字母，小写字母，数字、特殊字符四类中的三类，其中可输入的特殊字符为：~!@#$%&*()-_=|[{}];:,<.>/?
	Newpassword string `json:"newpassword"`
}

func (o ResetPasswordReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetPasswordReq struct{}"
	}

	return strings.Join([]string{"ResetPasswordReq", string(data)}, " ")
}
