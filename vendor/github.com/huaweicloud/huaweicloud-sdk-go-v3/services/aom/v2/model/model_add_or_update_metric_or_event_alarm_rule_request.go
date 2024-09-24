package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddOrUpdateMetricOrEventAlarmRuleRequest Request Object
type AddOrUpdateMetricOrEventAlarmRuleRequest struct {

	// 告警规则id。 - 新增告警时，填写\"add-alarm-action\" - 更新告警时，填写“update-alarm-action”
	ActionId string `json:"action_id"`

	// 企业项目id。 - 查询单个企业项目下实例，填写企业项目id。  - 查询所有企业项目下实例，填写“all_granted_eps”。
	EnterpriseProjectId *string `json:"Enterprise-Project-Id,omitempty"`

	Body *AddOrUpdateAlarmRuleV4RequestBody `json:"body,omitempty"`
}

func (o AddOrUpdateMetricOrEventAlarmRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddOrUpdateMetricOrEventAlarmRuleRequest struct{}"
	}

	return strings.Join([]string{"AddOrUpdateMetricOrEventAlarmRuleRequest", string(data)}, " ")
}
