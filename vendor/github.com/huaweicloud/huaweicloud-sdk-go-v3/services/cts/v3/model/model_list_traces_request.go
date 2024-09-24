package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListTracesRequest Request Object
type ListTracesRequest struct {

	// 标识审计事件类型。 目前支持管理类事件（system）和数据类事件（data）。 默认值为\"system\"。
	TraceType ListTracesRequestTraceType `json:"trace_type"`

	// 标示查询事件列表中限定返回的事件条数。不传时默认10条，最大值200条。
	Limit *int32 `json:"limit,omitempty"`

	// 标识查询事件列表的起始时间戳（timestamp，为标准UTC时间，毫秒级，13位数字，不包括传入时间）默认为上一小时的时间戳。查询条件from与to配套使用。
	From *int64 `json:"from,omitempty"`

	// 取值为响应中中marker的值，用于标识查询事件的起始时间（自此条事件的记录时间起，向更早时间查询）。 可以与“from”、“to”结合使用。 最终的查询条件取两组时间条件的交集。
	Next *string `json:"next,omitempty"`

	// 标识查询事件列表的结束时间戳（timestamp，为标准UTC时间，毫秒级，13位数字，不包括传入时间）默认为当前时间戳。查询条件to与from配套使用。
	To *int64 `json:"to,omitempty"`

	// 当\"trace_type\"字段值为\"system\"时，该字段值默认为\"system\"。 当\"trace_type\"字段值为\"data\"时，该字段值可以传入数据类追踪器名称，达到筛选某个数据类追踪器下的数据事件目的。
	TrackerName *string `json:"tracker_name,omitempty"`

	// 标识查询事件列表对应的云服务类型。必须为已对接CTS的云服务的英文缩写，且服务类型一般为大写字母。 当\"trace_type\"字段值为\"system\"时，该字段筛选有效\"。 已对接的云服务列表参见《云审计服务用户指南》“支持的服务”章节。
	ServiceType *string `json:"service_type,omitempty"`

	// 标识特定用户名称，用以查询该用户下的所有事件。 当\"trace_type\"字段值为\"system\"时，该字段筛选有效\"。
	User *string `json:"user,omitempty"`

	// 标示查询事件列表对应的云服务资源ID。 当\"trace_type\"字段值为\"system\"时，该字段筛选有效\"。
	ResourceId *string `json:"resource_id,omitempty"`

	// 标示查询事件列表对应的的资源名称。 当\"trace_type\"字段值为\"system\"时，该字段筛选有效\"。 说明：该字段可能包含大写字母。
	ResourceName *string `json:"resource_name,omitempty"`

	// 标示查询事件列表对应的资源类型。 当\"trace_type\"字段值为\"system\"时，该字段筛选有效\"。
	ResourceType *string `json:"resource_type,omitempty"`

	// 标示某一条事件的事件ID。当传入这个查询条件时，其他查询条件自动不生效。 当\"trace_type\"字段值为\"system\"时，该字段筛选有效\"。
	TraceId *string `json:"trace_id,omitempty"`

	// 标示查询事件列表对应的事件名称。 当\"trace_type\"字段值为\"system\"时，该字段筛选有效\"。 说明：该字段可能包含大写字母。
	TraceName *string `json:"trace_name,omitempty"`

	// 标示查询事件列表对应的事件等级目前有三种：正常(normal), 警告(warning),事故(incident)。 当\"trace_type\"字段值为\"system\"时，该字段筛选有效\"。
	TraceRating *ListTracesRequestTraceRating `json:"trace_rating,omitempty"`

	// 标示查询事件列表对应的访问密钥ID。包含临时访问凭证和永久访问密钥。
	AccessKeyId *string `json:"access_key_id,omitempty"`

	// 标示查询事件列表对应的企业项目ID。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ListTracesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTracesRequest struct{}"
	}

	return strings.Join([]string{"ListTracesRequest", string(data)}, " ")
}

type ListTracesRequestTraceType struct {
	value string
}

type ListTracesRequestTraceTypeEnum struct {
	SYSTEM ListTracesRequestTraceType
	DATA   ListTracesRequestTraceType
}

func GetListTracesRequestTraceTypeEnum() ListTracesRequestTraceTypeEnum {
	return ListTracesRequestTraceTypeEnum{
		SYSTEM: ListTracesRequestTraceType{
			value: "system",
		},
		DATA: ListTracesRequestTraceType{
			value: "data",
		},
	}
}

func (c ListTracesRequestTraceType) Value() string {
	return c.value
}

func (c ListTracesRequestTraceType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListTracesRequestTraceType) UnmarshalJSON(b []byte) error {
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

type ListTracesRequestTraceRating struct {
	value string
}

type ListTracesRequestTraceRatingEnum struct {
	NORMAL   ListTracesRequestTraceRating
	WARNING  ListTracesRequestTraceRating
	INCIDENT ListTracesRequestTraceRating
}

func GetListTracesRequestTraceRatingEnum() ListTracesRequestTraceRatingEnum {
	return ListTracesRequestTraceRatingEnum{
		NORMAL: ListTracesRequestTraceRating{
			value: "normal",
		},
		WARNING: ListTracesRequestTraceRating{
			value: "warning",
		},
		INCIDENT: ListTracesRequestTraceRating{
			value: "incident",
		},
	}
}

func (c ListTracesRequestTraceRating) Value() string {
	return c.value
}

func (c ListTracesRequestTraceRating) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListTracesRequestTraceRating) UnmarshalJSON(b []byte) error {
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
