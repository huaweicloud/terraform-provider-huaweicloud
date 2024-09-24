package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RunImageSynchronizeRequest Request Object
type RunImageSynchronizeRequest struct {

	// Region ID
	Region *string `json:"region,omitempty"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *RunImageSynchronizeRequestInfo `json:"body,omitempty"`
}

func (o RunImageSynchronizeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RunImageSynchronizeRequest struct{}"
	}

	return strings.Join([]string{"RunImageSynchronizeRequest", string(data)}, " ")
}
