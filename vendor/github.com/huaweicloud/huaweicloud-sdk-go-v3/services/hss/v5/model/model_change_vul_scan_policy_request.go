package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeVulScanPolicyRequest Request Object
type ChangeVulScanPolicyRequest struct {

	// 企业租户ID，查询所有企业项目时填写：all_granted_eps
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
