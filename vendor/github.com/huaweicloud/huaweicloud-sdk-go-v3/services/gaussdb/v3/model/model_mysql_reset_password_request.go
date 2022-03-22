package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlResetPasswordRequest struct {
	// 数据库密码。取值范围：至少包含以下字符的三种：大小写字母、数字和特殊符号~!@#$%^*-_=+?,()&，长度8~32个字符。建议您输入高强度密码，以提高安全性，防止出现密码被暴力破解等安全风险。如果您输入弱密码，系统会自动判定密码非法。

	Password string `json:"password"`
}

func (o MysqlResetPasswordRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlResetPasswordRequest struct{}"
	}

	return strings.Join([]string{"MysqlResetPasswordRequest", string(data)}, " ")
}
