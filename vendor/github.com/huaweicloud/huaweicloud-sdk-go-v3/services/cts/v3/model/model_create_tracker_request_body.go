package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type CreateTrackerRequestBody struct {

	// 标识追踪器类型。 目前支持系统追踪器类型有管理类追踪器(system)和数据类追踪器(data)。 数据类追踪器和管理类追踪器共同参数有：is_lts_enabled, obs_info; 管理类追踪器参数：is_support_trace_files_encryption, kms_id, is_support_validate, is_support_validate; 数据类追踪器参数：tracker_name, data_bucket。
	TrackerType CreateTrackerRequestBodyTrackerType `json:"tracker_type"`

	// 标识追踪器名称。 当\"tracker_type\"参数值为\"system\"时该参数为默认值\"system\"。 当\"tracker_type\"参数值为\"data\"时该参数需要指定追踪器名称\"。
	TrackerName string `json:"tracker_name"`

	// 云服务委托名称。 参数值为\"cts_admin_trust\"时，创建追踪器会自动创建一个云服务委托：cts_admin_trust。
	AgencyName *CreateTrackerRequestBodyAgencyName `json:"agency_name,omitempty"`

	// 是否应用到我的组织。 只针对管理类追踪器。设置为true时，ORG组织下所有成员当前区域的审计日志会转储到该追踪器配置的OBS桶或者LTS日志流，但是事件列表界面不支持查看其它组织成员的审计日志。
	IsOrganizationTracker *bool `json:"is_organization_tracker,omitempty"`

	ManagementEventSelector *ManagementEventSelector `json:"management_event_selector,omitempty"`

	// 是否打开事件分析。
	IsLtsEnabled *bool `json:"is_lts_enabled,omitempty"`

	ObsInfo *TrackerObsInfo `json:"obs_info,omitempty"`

	// 事件文件转储加密功能开关。 当\"tracker_type\"参数值为\"system\"时该参数值有效。 该参数必须与kms_id参数同时使用。
	IsSupportTraceFilesEncryption *bool `json:"is_support_trace_files_encryption,omitempty"`

	// 事件文件转储加密所采用的秘钥id（从KMS获取）。 当\"tracker_type\"参数值为\"system\"时该参数值有效。 当\"is_support_trace_files_encryption\"参数值为“是”时，此参数为必选项。
	KmsId *string `json:"kms_id,omitempty"`

	// 事件文件转储时是否打开事件文件校验。 当\"tracker_type\"参数值为\"system\"时该参数值有效。
	IsSupportValidate *bool `json:"is_support_validate,omitempty"`

	DataBucket *DataBucket `json:"data_bucket,omitempty"`
}

func (o CreateTrackerRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrackerRequestBody struct{}"
	}

	return strings.Join([]string{"CreateTrackerRequestBody", string(data)}, " ")
}

type CreateTrackerRequestBodyTrackerType struct {
	value string
}

type CreateTrackerRequestBodyTrackerTypeEnum struct {
	SYSTEM CreateTrackerRequestBodyTrackerType
	DATA   CreateTrackerRequestBodyTrackerType
}

func GetCreateTrackerRequestBodyTrackerTypeEnum() CreateTrackerRequestBodyTrackerTypeEnum {
	return CreateTrackerRequestBodyTrackerTypeEnum{
		SYSTEM: CreateTrackerRequestBodyTrackerType{
			value: "system",
		},
		DATA: CreateTrackerRequestBodyTrackerType{
			value: "data",
		},
	}
}

func (c CreateTrackerRequestBodyTrackerType) Value() string {
	return c.value
}

func (c CreateTrackerRequestBodyTrackerType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateTrackerRequestBodyTrackerType) UnmarshalJSON(b []byte) error {
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

type CreateTrackerRequestBodyAgencyName struct {
	value string
}

type CreateTrackerRequestBodyAgencyNameEnum struct {
	CTS_ADMIN_TRUST CreateTrackerRequestBodyAgencyName
}

func GetCreateTrackerRequestBodyAgencyNameEnum() CreateTrackerRequestBodyAgencyNameEnum {
	return CreateTrackerRequestBodyAgencyNameEnum{
		CTS_ADMIN_TRUST: CreateTrackerRequestBodyAgencyName{
			value: "cts_admin_trust",
		},
	}
}

func (c CreateTrackerRequestBodyAgencyName) Value() string {
	return c.value
}

func (c CreateTrackerRequestBodyAgencyName) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateTrackerRequestBodyAgencyName) UnmarshalJSON(b []byte) error {
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
