package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type RecordRule struct {

	// 规则id，由服务端返回。创建或修改的时候不携带
	Id *string `json:"id,omitempty"`

	// 直播推流域名
	PublishDomain string `json:"publish_domain"`

	// 应用名，如果任意应用填写*。录制规则匹配的时候，优先精确app匹配，如果匹配不到，则匹配*
	App string `json:"app"`

	// 录制的流名，如果任意流名则填写*。录制规则匹配的时候，优先精确stream匹配，如果匹配不到，则匹配*
	Stream string `json:"stream"`

	// 录制类型，包括：CONTINUOUS_RECORD，COMMAND_RECORD，PLAN_RECORD, ON_DEMAND_RECORD。默认CONTINUOUS_RECORD。 - CONTINUOUS_RECORD: 持续录制，在该规则类型配置后，只要有流到推送到录制系统，就触发录制。 - COMMAND_RECORD: 命令录制，在该规则类型配置后，在流推送到录制系统后，租户需要通过命令控制该流的录制开始和结束。命令控制的接口参考/v1/{project_id}/record/control - PLAN_RECORD: 计划录制，在该规则类型配置后，推的流如果在计划录制的时间区间则触发录制。 - ON_DEMAND_RECORD: 按需录制，在该规则类型配置后，录制系统收到推流后，需要调用租户提供的接口查询录制规则，并根据规则录制。租户提供的接口定义参考：/customer-record-ondemand-url
	RecordType *RecordRuleRecordType `json:"record_type,omitempty"`

	DefaultRecordConfig *DefaultRecordConfig `json:"default_record_config,omitempty"`

	// 创建时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间。 在查询的时候返回
	CreateTime *string `json:"create_time,omitempty"`

	// 修改时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间。 在查询的时候返回
	UpdateTime *string `json:"update_time,omitempty"`
}

func (o RecordRule) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecordRule struct{}"
	}

	return strings.Join([]string{"RecordRule", string(data)}, " ")
}

type RecordRuleRecordType struct {
	value string
}

type RecordRuleRecordTypeEnum struct {
	CONTINUOUS_RECORD RecordRuleRecordType
	COMMAND_RECORD    RecordRuleRecordType
	PLAN_RECORD       RecordRuleRecordType
	ON_DEMAND_RECORD  RecordRuleRecordType
}

func GetRecordRuleRecordTypeEnum() RecordRuleRecordTypeEnum {
	return RecordRuleRecordTypeEnum{
		CONTINUOUS_RECORD: RecordRuleRecordType{
			value: "CONTINUOUS_RECORD",
		},
		COMMAND_RECORD: RecordRuleRecordType{
			value: "COMMAND_RECORD",
		},
		PLAN_RECORD: RecordRuleRecordType{
			value: "PLAN_RECORD",
		},
		ON_DEMAND_RECORD: RecordRuleRecordType{
			value: "ON_DEMAND_RECORD",
		},
	}
}

func (c RecordRuleRecordType) Value() string {
	return c.value
}

func (c RecordRuleRecordType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RecordRuleRecordType) UnmarshalJSON(b []byte) error {
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
