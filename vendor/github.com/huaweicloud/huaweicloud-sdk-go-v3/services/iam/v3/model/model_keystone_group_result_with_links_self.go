package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneGroupResultWithLinksSelf struct {

	// 用户组描述信息。
	Description string `json:"description"`

	// 用户组ID。
	Id string `json:"id"`

	// 用户组所属账号ID。
	DomainId string `json:"domain_id"`

	// 用户组名称。
	Name string `json:"name"`

	Links *LinksSelf `json:"links"`

	// 用户组创建时间。
	CreateTime int64 `json:"create_time"`
}

func (o KeystoneGroupResultWithLinksSelf) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneGroupResultWithLinksSelf struct{}"
	}

	return strings.Join([]string{"KeystoneGroupResultWithLinksSelf", string(data)}, " ")
}
