package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListRecordContentsRequest Request Object
type ListRecordContentsRequest struct {

	// 直播推流域名
	PublishDomain *string `json:"publish_domain,omitempty"`

	// 流应用名称
	App *string `json:"app,omitempty"`

	// 流名称
	Stream *string `json:"stream,omitempty"`

	// 录制类型，包括：CONTINUOUS_RECORD，COMMAND_RECORD。默认CONTINUOUS_RECORD。 - CONTINUOUS_RECORD：持续录制，在该规则类型配置后，只要有流到推送到录制系统，就触发录制。 - COMMAND_RECORD：命令录制，在该规则类型配置后，在流推送到录制系统后，租户需要通过命令控制该流的录制开始和结束。
	RecordType *ListRecordContentsRequestRecordType `json:"record_type,omitempty"`

	// 开始时间,格式为：yyyy-mm-ddThh:mm:ssZ，UTC时间
	StartTime string `json:"start_time"`

	// 结束时间，格式为：yyyy-mm-ddThh:mm:ssZ，UTC时间
	EndTime *string `json:"end_time,omitempty"`

	// 分页编号，从0开始算
	Offset *int32 `json:"offset,omitempty"`

	// 每页记录数，取值范围[1,100]
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListRecordContentsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRecordContentsRequest struct{}"
	}

	return strings.Join([]string{"ListRecordContentsRequest", string(data)}, " ")
}

type ListRecordContentsRequestRecordType struct {
	value string
}

type ListRecordContentsRequestRecordTypeEnum struct {
	CONTINUOUS_RECORD ListRecordContentsRequestRecordType
	COMMAND_RECORD    ListRecordContentsRequestRecordType
}

func GetListRecordContentsRequestRecordTypeEnum() ListRecordContentsRequestRecordTypeEnum {
	return ListRecordContentsRequestRecordTypeEnum{
		CONTINUOUS_RECORD: ListRecordContentsRequestRecordType{
			value: "CONTINUOUS_RECORD",
		},
		COMMAND_RECORD: ListRecordContentsRequestRecordType{
			value: "COMMAND_RECORD",
		},
	}
}

func (c ListRecordContentsRequestRecordType) Value() string {
	return c.value
}

func (c ListRecordContentsRequestRecordType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListRecordContentsRequestRecordType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
