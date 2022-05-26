package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type RecordRuleRequest struct {

	// 直播推流域名
	PublishDomain string `json:"publish_domain"`

	// 应用名，如需匹配任意应用则填写*。录制规则匹配的时候，优先精确app匹配，如果匹配不到，则匹配*
	App string `json:"app"`

	// 录制的流名，如需匹配任流名则填写*。录制规则匹配的时候，优先精确stream匹配，如果匹配不到，则匹配*
	Stream string `json:"stream"`

	// 录制类型，包括：CONTINUOUS_RECORD，COMMAND_RECORD，PLAN_RECORD，ON_DEMAND_RECORD。默认CONTINUOUS_RECORD。 - CONTINUOUS_RECORD：持续录制，在该规则类型配置后，只要有流到推送到录制系统，就触发录制。 - COMMAND_RECORD：命令录制，在该规则类型配置后，在流推送到录制系统后，租户需要通过命令控制该流的录制开始和结束。 - PLAN_RECORD：计划录制，在该规则类型配置后，推的流如果在计划录制的时间区间则触发录制。 - ON_DEMAND_RECORD：按需录制，在该规则类型配置后，录制系统收到推流后，需要调用租户提供的接口查询录制规则，并根据规则录制。
	RecordType *RecordRuleRequestRecordType `json:"record_type,omitempty"`

	DefaultRecordConfig *DefaultRecordConfig `json:"default_record_config"`
}

func (o RecordRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecordRuleRequest struct{}"
	}

	return strings.Join([]string{"RecordRuleRequest", string(data)}, " ")
}

type RecordRuleRequestRecordType struct {
	value string
}

type RecordRuleRequestRecordTypeEnum struct {
	CONTINUOUS_RECORD RecordRuleRequestRecordType
	COMMAND_RECORD    RecordRuleRequestRecordType
	PLAN_RECORD       RecordRuleRequestRecordType
	ON_DEMAND_RECORD  RecordRuleRequestRecordType
}

func GetRecordRuleRequestRecordTypeEnum() RecordRuleRequestRecordTypeEnum {
	return RecordRuleRequestRecordTypeEnum{
		CONTINUOUS_RECORD: RecordRuleRequestRecordType{
			value: "CONTINUOUS_RECORD",
		},
		COMMAND_RECORD: RecordRuleRequestRecordType{
			value: "COMMAND_RECORD",
		},
		PLAN_RECORD: RecordRuleRequestRecordType{
			value: "PLAN_RECORD",
		},
		ON_DEMAND_RECORD: RecordRuleRequestRecordType{
			value: "ON_DEMAND_RECORD",
		},
	}
}

func (c RecordRuleRequestRecordType) Value() string {
	return c.value
}

func (c RecordRuleRequestRecordType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RecordRuleRequestRecordType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
