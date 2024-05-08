package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type PromInstanceEpsModel struct {

	// Prometheus实例名称 名称不能以下划线或中划线开头结尾，只含有中文、英文、数字、下划线、中划线、长度1-100。
	PromName string `json:"prom_name"`

	// Prometheus实例id。
	PromId *string `json:"prom_id,omitempty"`

	// Prometheus实例类型（暂时不支持VPC、KUBERNETES）。
	PromType PromInstanceEpsModelPromType `json:"prom_type"`

	// Prometheus实例版本号。
	PromVersion *string `json:"prom_version,omitempty"`

	// Prometheus实例创建时间戳。
	PromCreateTimestamp *int64 `json:"prom_create_timestamp,omitempty"`

	// Prometheus实例更新时间戳。
	PromUpdateTimestamp *int64 `json:"prom_update_timestamp,omitempty"`

	// Prometheus实例状态。
	PromStatus *PromInstanceEpsModelPromStatus `json:"prom_status,omitempty"`

	// Prometheus实例所属的企业项目。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// Prometheus实例所属projectId。
	ProjectId *string `json:"project_id,omitempty"`

	// 删除标记。
	IsDeletedTag *int64 `json:"is_deleted_tag,omitempty"`

	// 删除时间。
	DeletedTime *int64 `json:"deleted_time,omitempty"`

	PromSpecConfig *PromConfigModel `json:"prom_spec_config,omitempty"`

	// Prometheus实例所属CCE特殊配置。
	CceSpecConfig *string `json:"cce_spec_config,omitempty"`
}

func (o PromInstanceEpsModel) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PromInstanceEpsModel struct{}"
	}

	return strings.Join([]string{"PromInstanceEpsModel", string(data)}, " ")
}

type PromInstanceEpsModelPromType struct {
	value string
}

type PromInstanceEpsModelPromTypeEnum struct {
	DEFAULT        PromInstanceEpsModelPromType
	ECS            PromInstanceEpsModelPromType
	VPC            PromInstanceEpsModelPromType
	CCE            PromInstanceEpsModelPromType
	REMOTE_WRITE   PromInstanceEpsModelPromType
	KUBERNETES     PromInstanceEpsModelPromType
	CLOUD_SERVICE  PromInstanceEpsModelPromType
	ACROSS_ACCOUNT PromInstanceEpsModelPromType
}

func GetPromInstanceEpsModelPromTypeEnum() PromInstanceEpsModelPromTypeEnum {
	return PromInstanceEpsModelPromTypeEnum{
		DEFAULT: PromInstanceEpsModelPromType{
			value: "default",
		},
		ECS: PromInstanceEpsModelPromType{
			value: "ECS",
		},
		VPC: PromInstanceEpsModelPromType{
			value: "VPC",
		},
		CCE: PromInstanceEpsModelPromType{
			value: "CCE",
		},
		REMOTE_WRITE: PromInstanceEpsModelPromType{
			value: "REMOTE_WRITE",
		},
		KUBERNETES: PromInstanceEpsModelPromType{
			value: "KUBERNETES",
		},
		CLOUD_SERVICE: PromInstanceEpsModelPromType{
			value: "CLOUD_SERVICE",
		},
		ACROSS_ACCOUNT: PromInstanceEpsModelPromType{
			value: "ACROSS_ACCOUNT",
		},
	}
}

func (c PromInstanceEpsModelPromType) Value() string {
	return c.value
}

func (c PromInstanceEpsModelPromType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PromInstanceEpsModelPromType) UnmarshalJSON(b []byte) error {
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

type PromInstanceEpsModelPromStatus struct {
	value string
}

type PromInstanceEpsModelPromStatusEnum struct {
	DELETED PromInstanceEpsModelPromStatus
	NORMAL  PromInstanceEpsModelPromStatus
	ALL     PromInstanceEpsModelPromStatus
}

func GetPromInstanceEpsModelPromStatusEnum() PromInstanceEpsModelPromStatusEnum {
	return PromInstanceEpsModelPromStatusEnum{
		DELETED: PromInstanceEpsModelPromStatus{
			value: "DELETED",
		},
		NORMAL: PromInstanceEpsModelPromStatus{
			value: "NORMAL",
		},
		ALL: PromInstanceEpsModelPromStatus{
			value: "ALL",
		},
	}
}

func (c PromInstanceEpsModelPromStatus) Value() string {
	return c.value
}

func (c PromInstanceEpsModelPromStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PromInstanceEpsModelPromStatus) UnmarshalJSON(b []byte) error {
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
