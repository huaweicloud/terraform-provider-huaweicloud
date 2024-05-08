package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopProtectionRequest Request Object
type StopProtectionRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *CloseProtectionInfoRequestInfo `json:"body,omitempty"`
}

func (o StopProtectionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopProtectionRequest struct{}"
	}

	return strings.Join([]string{"StopProtectionRequest", string(data)}, " ")
}
