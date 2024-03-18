package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type CreateSmartConnectTaskReq struct {

	// SmartConnect任务名称。
	TaskName *string `json:"task_name,omitempty"`

	// 是否稍后再启动任务。如需要创建任务后立即启动，请填false；如希望稍后在任务列表中手动开启任务，请填true。
	StartLater *bool `json:"start_later,omitempty"`

	// SmartConnect任务配置的Topic。
	Topics *string `json:"topics,omitempty"`

	// SmartConnect任务配置的Topic正则表达式。
	TopicsRegex *string `json:"topics_regex,omitempty"`

	// SmartConnect任务的源端类型。
	SourceType *CreateSmartConnectTaskReqSourceType `json:"source_type,omitempty"`

	SourceTask *SmartConnectTaskReqSourceConfig `json:"source_task,omitempty"`

	// SmartConnect任务的目标端类型。
	SinkType *CreateSmartConnectTaskReqSinkType `json:"sink_type,omitempty"`

	SinkTask *SmartConnectTaskReqSinkConfig `json:"sink_task,omitempty"`
}

func (o CreateSmartConnectTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSmartConnectTaskReq struct{}"
	}

	return strings.Join([]string{"CreateSmartConnectTaskReq", string(data)}, " ")
}

type CreateSmartConnectTaskReqSourceType struct {
	value string
}

type CreateSmartConnectTaskReqSourceTypeEnum struct {
	REDIS_REPLICATOR_SOURCE CreateSmartConnectTaskReqSourceType
	KAFKA_REPLICATOR_SOURCE CreateSmartConnectTaskReqSourceType
	NONE                    CreateSmartConnectTaskReqSourceType
}

func GetCreateSmartConnectTaskReqSourceTypeEnum() CreateSmartConnectTaskReqSourceTypeEnum {
	return CreateSmartConnectTaskReqSourceTypeEnum{
		REDIS_REPLICATOR_SOURCE: CreateSmartConnectTaskReqSourceType{
			value: "REDIS_REPLICATOR_SOURCE",
		},
		KAFKA_REPLICATOR_SOURCE: CreateSmartConnectTaskReqSourceType{
			value: "KAFKA_REPLICATOR_SOURCE",
		},
		NONE: CreateSmartConnectTaskReqSourceType{
			value: "NONE",
		},
	}
}

func (c CreateSmartConnectTaskReqSourceType) Value() string {
	return c.value
}

func (c CreateSmartConnectTaskReqSourceType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateSmartConnectTaskReqSourceType) UnmarshalJSON(b []byte) error {
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

type CreateSmartConnectTaskReqSinkType struct {
	value string
}

type CreateSmartConnectTaskReqSinkTypeEnum struct {
	REDIS_REPLICATOR_SINK CreateSmartConnectTaskReqSinkType
	OBS_SINK              CreateSmartConnectTaskReqSinkType
	NONE                  CreateSmartConnectTaskReqSinkType
}

func GetCreateSmartConnectTaskReqSinkTypeEnum() CreateSmartConnectTaskReqSinkTypeEnum {
	return CreateSmartConnectTaskReqSinkTypeEnum{
		REDIS_REPLICATOR_SINK: CreateSmartConnectTaskReqSinkType{
			value: "REDIS_REPLICATOR_SINK",
		},
		OBS_SINK: CreateSmartConnectTaskReqSinkType{
			value: "OBS_SINK",
		},
		NONE: CreateSmartConnectTaskReqSinkType{
			value: "NONE",
		},
	}
}

func (c CreateSmartConnectTaskReqSinkType) Value() string {
	return c.value
}

func (c CreateSmartConnectTaskReqSinkType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateSmartConnectTaskReqSinkType) UnmarshalJSON(b []byte) error {
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
