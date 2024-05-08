package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBackupPolicyInfoRequest Request Object
type ShowBackupPolicyInfoRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ShowBackupPolicyInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBackupPolicyInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowBackupPolicyInfoRequest", string(data)}, " ")
}
