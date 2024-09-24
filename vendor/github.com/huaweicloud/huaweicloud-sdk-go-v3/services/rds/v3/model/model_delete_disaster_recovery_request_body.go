package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDisasterRecoveryRequestBody 解除实例容灾关系请求体
type DeleteDisasterRecoveryRequestBody struct {

	// 解除目标的实例id
	TargetInstanceId string `json:"target_instance_id"`

	// 解除目标的project id
	TargetProjectId string `json:"target_project_id"`

	// 解除目标的region
	TargetRegion string `json:"target_region"`

	// 解除目标的数据浮动ip
	TargetIp string `json:"target_ip"`

	// 当前操作对象是否是主实例
	IsMaster bool `json:"is_master"`
}

func (o DeleteDisasterRecoveryRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDisasterRecoveryRequestBody struct{}"
	}

	return strings.Join([]string{"DeleteDisasterRecoveryRequestBody", string(data)}, " ")
}
