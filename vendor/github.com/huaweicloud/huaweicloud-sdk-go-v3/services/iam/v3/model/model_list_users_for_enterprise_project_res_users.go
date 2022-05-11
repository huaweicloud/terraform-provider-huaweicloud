package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ListUsersForEnterpriseProjectResUsers struct {

	// 授权用户所属账号ID。
	DomainId string `json:"domain_id"`

	// 授权用户ID。
	Id string `json:"id"`

	// 授权用户名。
	Name string `json:"name"`

	// 授权用户是否启用，true表示启用，false表示停用，默认为true。
	Enabled bool `json:"enabled"`

	// 授权用户描述信息。
	Description string `json:"description"`

	// 授权用户的策略数。
	PolicyNum int32 `json:"policy_num"`

	// 用户最近与企业项目关联策略的时间（毫秒）。
	LastestPolicyTime int64 `json:"lastest_policy_time"`
}

func (o ListUsersForEnterpriseProjectResUsers) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUsersForEnterpriseProjectResUsers struct{}"
	}

	return strings.Join([]string{"ListUsersForEnterpriseProjectResUsers", string(data)}, " ")
}
