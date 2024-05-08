package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetWtpProtectionStatusInfoRequest Request Object
type SetWtpProtectionStatusInfoRequest struct {

	// 缺省值:application/json; charset=utf-8
	ContentType *string `json:"Content-Type,omitempty"`

	// Region Id
	Region string `json:"region"`

	// 企业项目ID
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
