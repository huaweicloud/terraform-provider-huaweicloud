package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ListGroupsForEnterpriseProjectResDetail struct {

	// 用户组创建时间。
	CreateTime int64 `json:"createTime"`

	// 用户组描述。
	Description string `json:"description"`

	// 租户Id。
	DomainId string `json:"domainId"`

	// 用户组Id。
	Id string `json:"id"`

	// 用户组名称。
	Name string `json:"name"`
}

func (o ListGroupsForEnterpriseProjectResDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListGroupsForEnterpriseProjectResDetail struct{}"
	}

	return strings.Join([]string{"ListGroupsForEnterpriseProjectResDetail", string(data)}, " ")
}
