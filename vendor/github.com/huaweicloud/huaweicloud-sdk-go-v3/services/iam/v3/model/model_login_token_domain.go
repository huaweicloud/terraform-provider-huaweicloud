package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type LoginTokenDomain struct {

	// 被委托方用户所属账号名称。
	Name *string `json:"name,omitempty"`

	// 被委托方用户所属账号ID。
	Id *string `json:"id,omitempty"`
}

func (o LoginTokenDomain) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LoginTokenDomain struct{}"
	}

	return strings.Join([]string{"LoginTokenDomain", string(data)}, " ")
}
