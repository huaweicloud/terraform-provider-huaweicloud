package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteAccountRequest Request Object
type DeleteAccountRequest struct {

	// 如果正在使用临时安全凭据，则此header是必需的，该值是临时安全凭据的安全令牌（会话令牌）。
	XSecurityToken *string `json:"X-Security-Token,omitempty"`

	// Region ID
	Region string `json:"region"`

	// 企业租户ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *DeleteAccountRequestInfo `json:"body,omitempty"`
}

func (o DeleteAccountRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAccountRequest struct{}"
	}

	return strings.Join([]string{"DeleteAccountRequest", string(data)}, " ")
}
