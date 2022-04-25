package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateSecurityGroupResponse struct {

	// 请求Id
	RequestId *string `json:"request_id,omitempty"`

	SecurityGroup  *SecurityGroupInfo `json:"security_group,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o CreateSecurityGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSecurityGroupResponse struct{}"
	}

	return strings.Join([]string{"CreateSecurityGroupResponse", string(data)}, " ")
}
