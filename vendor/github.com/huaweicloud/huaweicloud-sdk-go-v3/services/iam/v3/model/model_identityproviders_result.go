package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type IdentityprovidersResult struct {

	// 身份提供商类型。当前支持virtual_user_sso和iam_user_sso两种。当返回为空字符串或者null时，默认为缺省类型virtual_user_sso类型。
	SsoType string `json:"sso_type"`

	// 身份提供商ID。
	Id string `json:"id"`

	// 身份提供商描述信息。
	Description string `json:"description"`

	// 身份提供商是否启用，true为启用，false为停用，默认为false。
	Enabled bool `json:"enabled"`

	// 身份提供商的联邦用户ID列表。
	RemoteIds []string `json:"remote_ids"`

	Links *IdentityprovidersLinks `json:"links"`
}

func (o IdentityprovidersResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IdentityprovidersResult struct{}"
	}

	return strings.Join([]string{"IdentityprovidersResult", string(data)}, " ")
}
