package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type AssociatePolicyGroupRequest struct {

	// region id
	Region string `json:"region"`

	// 租户企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *AssociatePolicyGroupRequestInfo `json:"body,omitempty"`
}

func (o AssociatePolicyGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociatePolicyGroupRequest struct{}"
	}

	return strings.Join([]string{"AssociatePolicyGroupRequest", string(data)}, " ")
}
