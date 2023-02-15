package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type SetWtpProtectionStatusInfoRequest struct {

	// Region Id
	Region string `json:"region"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *SetWtpProtectionStatusRequestInfo `json:"body,omitempty"`
}

func (o SetWtpProtectionStatusInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetWtpProtectionStatusInfoRequest struct{}"
	}

	return strings.Join([]string{"SetWtpProtectionStatusInfoRequest", string(data)}, " ")
}
