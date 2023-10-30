package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateSecurityGroupRulesRequest Request Object
type BatchCreateSecurityGroupRulesRequest struct {

	// 安全组ID
	SecurityGroupId string `json:"security_group_id"`

	Body *BatchCreateSecurityGroupRulesRequestBody `json:"body,omitempty"`
}

func (o BatchCreateSecurityGroupRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateSecurityGroupRulesRequest struct{}"
	}

	return strings.Join([]string{"BatchCreateSecurityGroupRulesRequest", string(data)}, " ")
}
