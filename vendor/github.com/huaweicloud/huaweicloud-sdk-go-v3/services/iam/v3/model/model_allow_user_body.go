package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AllowUserBody 用户可以自主修改的属性。
type AllowUserBody struct {

	// 是否允许子用户自行管理AK，取值范围true或false。
	ManageAccesskey *bool `json:"manage_accesskey,omitempty"`

	// 是否允许子用户自己修改邮箱，取值范围true或false。
	ManageEmail *bool `json:"manage_email,omitempty"`

	// 是否允许子用户自己修改电话，取值范围true或false。
	ManageMobile *bool `json:"manage_mobile,omitempty"`

	// 是否允许子用户自己修改密码，取值范围true或false。
	ManagePassword *bool `json:"manage_password,omitempty"`
}

func (o AllowUserBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AllowUserBody struct{}"
	}

	return strings.Join([]string{"AllowUserBody", string(data)}, " ")
}
