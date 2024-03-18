package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SwitchSqlLimitRequest Request Object
type SwitchSqlLimitRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *SwitchSqlLimitControlReqV3 `json:"body,omitempty"`
}

func (o SwitchSqlLimitRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SwitchSqlLimitRequest struct{}"
	}

	return strings.Join([]string{"SwitchSqlLimitRequest", string(data)}, " ")
}
