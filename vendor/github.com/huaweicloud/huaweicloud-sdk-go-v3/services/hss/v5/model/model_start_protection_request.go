package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartProtectionRequest Request Object
type StartProtectionRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *ProtectionInfoRequestInfo `json:"body,omitempty"`
}

func (o StartProtectionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartProtectionRequest struct{}"
	}

	return strings.Join([]string{"StartProtectionRequest", string(data)}, " ")
}
