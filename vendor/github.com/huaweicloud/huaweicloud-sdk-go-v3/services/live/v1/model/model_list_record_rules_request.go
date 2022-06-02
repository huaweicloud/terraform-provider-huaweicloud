package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListRecordRulesRequest struct {

	// 直播推流域名
	PublishDomain *string `json:"publish_domain,omitempty"`

	// 流应用名称
	App *string `json:"app,omitempty"`

	// 流名称
	Stream *string `json:"stream,omitempty"`

	// 录制类型，如果不填写则查询所有录制类型，包括：CONTINUOUS_RECORD，COMMAND_RECORD，PLAN_RECORD，ON_DEMAND_RECORD。默认CONTINUOUS_RECORD。 - CONTINUOUS_RECORD：持续录制，在该规则类型配置后，只要有流到推送到录制系统，就触发录制。 - COMMAND_RECORD：命令录制，在该规则类型配置后，在流推送到录制系统后，租户需要通过命令控制该流的录制开始和结束。 - PLAN_RECORD：计划录制，在该规则类型配置后，推的流如果在计划录制的时间区间则触发录制。 - ON_DEMAND_RECORD：按需录制，在该规则类型配置后，录制系统收到推流后，需要调用租户提供的接口查询录制规则，并根据规则录制。
	RecordType *ListRecordRulesRequestRecordType `json:"record_type,omitempty"`

	// 偏移量，表示从此偏移量开始查询，offset大于等于0
	Offset *int32 `json:"offset,omitempty"`

	// 每页记录数，取值范围[1,100]，默认值10
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListRecordRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRecordRulesRequest struct{}"
	}

	return strings.Join([]string{"ListRecordRulesRequest", string(data)}, " ")
}

type ListRecordRulesRequestRecordType struct {
	value string
}

type ListRecordRulesRequestRecordTypeEnum struct {
	CONTINUOUS_RECORD ListRecordRulesRequestRecordType
	COMMAND_RECORD    ListRecordRulesRequestRecordType
	PLAN_RECORD       ListRecordRulesRequestRecordType
	ON_DEMAND_RECORD  ListRecordRulesRequestRecordType
}

func GetListRecordRulesRequestRecordTypeEnum() ListRecordRulesRequestRecordTypeEnum {
	return ListRecordRulesRequestRecordTypeEnum{
		CONTINUOUS_RECORD: ListRecordRulesRequestRecordType{
			value: "CONTINUOUS_RECORD",
		},
		COMMAND_RECORD: ListRecordRulesRequestRecordType{
			value: "COMMAND_RECORD",
		},
		PLAN_RECORD: ListRecordRulesRequestRecordType{
			value: "PLAN_RECORD",
		},
		ON_DEMAND_RECORD: ListRecordRulesRequestRecordType{
			value: "ON_DEMAND_RECORD",
		},
	}
}

func (c ListRecordRulesRequestRecordType) Value() string {
	return c.value
}

func (c ListRecordRulesRequestRecordType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListRecordRulesRequestRecordType) UnmarshalJSON(b []byte) error {
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
