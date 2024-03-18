package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSqlLimitRequest Request Object
type DeleteSqlLimitRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *DeleteSqlLimitRuleReqV3 `json:"body,omitempty"`
}

func (o DeleteSqlLimitRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSqlLimitRequest struct{}"
	}

	return strings.Join([]string{"DeleteSqlLimitRequest", string(data)}, " ")
}
