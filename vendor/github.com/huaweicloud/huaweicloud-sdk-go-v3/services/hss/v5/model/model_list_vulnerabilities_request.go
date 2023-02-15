package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListVulnerabilitiesRequest struct {

	// 企业租户ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 漏洞类型，包含如下：   -linux_vul : linux漏洞   -windows_vul : windows漏洞   -web_cms : Web-CMS漏洞
	Type *string `json:"type,omitempty"`

	// 漏洞ID
	VulId *string `json:"vul_id,omitempty"`

	// 漏洞名称
	VulName *string `json:"vul_name,omitempty"`

	// 每页显示个数
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置，必须为数字，取值范围为大于或等于0，默认0
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListVulnerabilitiesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListVulnerabilitiesRequest struct{}"
	}

	return strings.Join([]string{"ListVulnerabilitiesRequest", string(data)}, " ")
}
