package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeBlockedIpRequest Request Object
type ChangeBlockedIpRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *ChangeBlockedIpRequestInfo `json:"body,omitempty"`
}

func (o ChangeBlockedIpRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeBlockedIpRequest struct{}"
	}

	return strings.Join([]string{"ChangeBlockedIpRequest", string(data)}, " ")
}
