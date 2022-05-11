package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type IdentityproviderOption struct {

	// 身份提供商类型。当前支持virtual_user_sso和iam_user_sso两种，缺省配置默认为virtual_user_sso类型。
	SsoType *string `json:"sso_type,omitempty"`

	// 身份提供商描述信息。
	Description *string `json:"description,omitempty"`

	// 身份提供商是否启用，true为启用，false为停用，默认为false。
	Enabled *bool `json:"enabled,omitempty"`
}

func (o IdentityproviderOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IdentityproviderOption struct{}"
	}

	return strings.Join([]string{"IdentityproviderOption", string(data)}, " ")
}
