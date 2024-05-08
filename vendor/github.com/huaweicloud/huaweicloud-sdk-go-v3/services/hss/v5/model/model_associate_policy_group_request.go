package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AssociatePolicyGroupRequest Request Object
type AssociatePolicyGroupRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 缺省值:application/json; charset=utf-8
	ContentType *string `json:"Content-Type,omitempty"`

	Body *AssociatePolicyGroupRequestInfo `json:"body,omitempty"`
}

func (o AssociatePolicyGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociatePolicyGroupRequest struct{}"
	}

	return strings.Join([]string{"AssociatePolicyGroupRequest", string(data)}, " ")
}
