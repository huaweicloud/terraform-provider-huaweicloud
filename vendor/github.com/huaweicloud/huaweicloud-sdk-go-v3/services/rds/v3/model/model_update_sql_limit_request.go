package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateSqlLimitRequest Request Object
type UpdateSqlLimitRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *UpdateSqlLimitRuleReqV3 `json:"body,omitempty"`
}

func (o UpdateSqlLimitRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSqlLimitRequest struct{}"
	}

	return strings.Join([]string{"UpdateSqlLimitRequest", string(data)}, " ")
}
