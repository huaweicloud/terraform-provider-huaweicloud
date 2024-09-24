package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListProtectionPolicyRequest Request Object
type ListProtectionPolicyRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`

	// 每页显示个数
	Limit *int32 `json:"limit,omitempty"`

	// 防护策略名称
	PolicyName *string `json:"policy_name,omitempty"`

	// 防护策略id
	ProtectPolicyId *string `json:"protect_policy_id,omitempty"`

	// 策略支持的操作系统，包含如下：   - Windows : Windows系统   - Linux : Linux系统
	OperatingSystem *string `json:"operating_system,omitempty"`
}

func (o ListProtectionPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProtectionPolicyRequest struct{}"
	}

	return strings.Join([]string{"ListProtectionPolicyRequest", string(data)}, " ")
}
