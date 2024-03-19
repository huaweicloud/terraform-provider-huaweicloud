package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateSqlLimitRequest Request Object
type CreateSqlLimitRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *CreateSqlLimitRuleReqV3 `json:"body,omitempty"`
}

func (o CreateSqlLimitRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSqlLimitRequest struct{}"
	}

	return strings.Join([]string{"CreateSqlLimitRequest", string(data)}, " ")
}
