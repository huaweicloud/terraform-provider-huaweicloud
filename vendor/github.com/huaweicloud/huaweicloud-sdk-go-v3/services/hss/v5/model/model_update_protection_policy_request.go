package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateProtectionPolicyRequest Request Object
type UpdateProtectionPolicyRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *UpdateProtectionPolicyInfoRequestInfo `json:"body,omitempty"`
}

func (o UpdateProtectionPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProtectionPolicyRequest struct{}"
	}

	return strings.Join([]string{"UpdateProtectionPolicyRequest", string(data)}, " ")
}
