package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateIdentityproviderOption struct {

	// 身份提供商描述信息。
	Description *string `json:"description,omitempty"`

	// 身份提供商是否启用，true为启用，false为停用，默认为false。
	Enabled *bool `json:"enabled,omitempty"`
}

func (o UpdateIdentityproviderOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateIdentityproviderOption struct{}"
	}

	return strings.Join([]string{"UpdateIdentityproviderOption", string(data)}, " ")
}
