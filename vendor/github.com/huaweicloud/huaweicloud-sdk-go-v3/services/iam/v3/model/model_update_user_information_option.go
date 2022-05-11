package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateUserInformationOption struct {

	// IAM用户的新邮箱，符合邮箱格式，长度小于等于255字符。
	Email *string `json:"email,omitempty"`

	// IAM用户的国家码+新手机号，手机号为纯数字，长度小于等于32字符。
	Mobile *string `json:"mobile,omitempty"`
}

func (o UpdateUserInformationOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateUserInformationOption struct{}"
	}

	return strings.Join([]string{"UpdateUserInformationOption", string(data)}, " ")
}
