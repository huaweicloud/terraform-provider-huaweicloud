package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeIsolatedFileRequest Request Object
type ChangeIsolatedFileRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *ChangeIsolatedFileRequestInfo `json:"body,omitempty"`
}

func (o ChangeIsolatedFileRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeIsolatedFileRequest struct{}"
	}

	return strings.Join([]string{"ChangeIsolatedFileRequest", string(data)}, " ")
}
