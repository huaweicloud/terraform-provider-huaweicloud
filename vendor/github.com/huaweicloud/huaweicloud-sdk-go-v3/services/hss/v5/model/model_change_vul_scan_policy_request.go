package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeVulScanPolicyRequest Request Object
type ChangeVulScanPolicyRequest struct {

	// 企业租户ID，注：修改漏洞扫描策略将影响租户账号下所有主机的漏洞扫描行为，因此开通了多企业项目的用户，该参数须填写“all_granted_eps”才能执行漏洞策略修改。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *ChangeVulScanPolicyRequestInfo `json:"body,omitempty"`
}

func (o ChangeVulScanPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeVulScanPolicyRequest struct{}"
	}

	return strings.Join([]string{"ChangeVulScanPolicyRequest", string(data)}, " ")
}
