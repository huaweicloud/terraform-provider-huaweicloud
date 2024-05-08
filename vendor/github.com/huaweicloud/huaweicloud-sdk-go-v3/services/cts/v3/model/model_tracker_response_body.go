package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type TrackerResponseBody struct {

	// 追踪器唯一标识。
	Id *string `json:"id,omitempty"`

	// 追踪器创建时间戳。
	CreateTime *int64 `json:"create_time,omitempty"`

	// 事件文件转储加密所采用的秘钥id（从KMS获取）。 当\"tracker_type\"参数值为\"system\"和\"is_support_trace_files_encryption\"参数值为“是”时，此参数为必选项。
	KmsId *string `json:"kms_id,omitempty"`

	// 是否打开事件文件校验。
	IsSupportValidate *bool `json:"is_support_validate,omitempty"`

	// 是否应用到我的组织。 只针对管理类追踪器。设置为true时，ORG组织下所有成员当前区域的审计日志会转储到该追踪器配置的OBS桶或者LTS日志流，但是事件列表界面不支持查看其它组织成员的审计日志。
	IsOrganizationTracker *bool `json:"is_organization_tracker,omitempty"`

	// 云服务委托名称。
	AgencyName *TrackerResponseBodyAgencyName `json:"agency_name,omitempty"`

	ManagementEventSelector *ManagementEventSelector `json:"management_event_selector,omitempty"`

	Lts *Lts `json:"lts,omitempty"`

	// 标识追踪器类型。 目前支持系统追踪器类型有管理类追踪器（system）和数据类追踪器（data）。
	TrackerType *TrackerResponseBodyTrackerType `json:"tracker_type,omitempty"`

	// 账号ID，参见《云审计服务API参考》“获取账号ID和项目ID”章节。
	DomainId *string `json:"domain_id,omitempty"`

	// 项目ID。
	ProjectId *string `json:"project_id,omitempty"`

	// 标识追踪器名称，当前版本默认为“system”。
	TrackerName *string `json:"tracker_name,omitempty"`

	// 标识追踪器状态，包括正常（enabled），停止（disabled）和异常（error）三种状态，状态为异常时需通过明细（detail）字段说明错误来源。
	Status *TrackerResponseBodyStatus `json:"status,omitempty"`

	// 该参数仅在追踪器状态异常时返回，用于标识追踪器异常的原因，包括桶策略异常（bucketPolicyError），桶不存在（noBucket）和欠费或冻结（arrears）三种原因。
	Detail *string `json:"detail,omitempty"`

	// 事件文件转储加密功能开关。 该参数必须与kms_id参数同时使用。 当前环境仅\"tracker_type\"参数值为\"system\"时支持该功能。
	IsSupportTraceFilesEncryption *bool `json:"is_support_trace_files_encryption,omitempty"`

	// LTS服务日志组的ID。
	GroupId *string `json:"group_id,omitempty"`

	// LTS服务日志流的ID。
	StreamId *string `json:"stream_id,omitempty"`

	ObsInfo *ObsInfo `json:"obs_info,omitempty"`

	DataBucket *DataBucketQuery `json:"data_bucket,omitempty"`
}

func (o TrackerResponseBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TrackerResponseBody struct{}"
	}

	return strings.Join([]string{"TrackerResponseBody", string(data)}, " ")
}

type TrackerResponseBodyAgencyName struct {
	value string
}

type TrackerResponseBodyAgencyNameEnum struct {
	CTS_ADMIN_TRUST TrackerResponseBodyAgencyName
}

func GetTrackerResponseBodyAgencyNameEnum() TrackerResponseBodyAgencyNameEnum {
	return TrackerResponseBodyAgencyNameEnum{
		CTS_ADMIN_TRUST: TrackerResponseBodyAgencyName{
			value: "cts_admin_trust",
		},
	}
}

func (c TrackerResponseBodyAgencyName) Value() string {
	return c.value
}

func (c TrackerResponseBodyAgencyName) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TrackerResponseBodyAgencyName) UnmarshalJSON(b []byte) error {
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

type TrackerResponseBodyTrackerType struct {
	value string
}

type TrackerResponseBodyTrackerTypeEnum struct {
	SYSTEM TrackerResponseBodyTrackerType
	DATA   TrackerResponseBodyTrackerType
}

func GetTrackerResponseBodyTrackerTypeEnum() TrackerResponseBodyTrackerTypeEnum {
	return TrackerResponseBodyTrackerTypeEnum{
		SYSTEM: TrackerResponseBodyTrackerType{
			value: "system",
		},
		DATA: TrackerResponseBodyTrackerType{
			value: "data",
		},
	}
}

func (c TrackerResponseBodyTrackerType) Value() string {
	return c.value
}

func (c TrackerResponseBodyTrackerType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TrackerResponseBodyTrackerType) UnmarshalJSON(b []byte) error {
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

type TrackerResponseBodyStatus struct {
	value string
}

type TrackerResponseBodyStatusEnum struct {
	ENABLED  TrackerResponseBodyStatus
	DISABLED TrackerResponseBodyStatus
}

func GetTrackerResponseBodyStatusEnum() TrackerResponseBodyStatusEnum {
	return TrackerResponseBodyStatusEnum{
		ENABLED: TrackerResponseBodyStatus{
			value: "enabled",
		},
		DISABLED: TrackerResponseBodyStatus{
			value: "disabled",
		},
	}
}

func (c TrackerResponseBodyStatus) Value() string {
	return c.value
}

func (c TrackerResponseBodyStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TrackerResponseBodyStatus) UnmarshalJSON(b []byte) error {
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
