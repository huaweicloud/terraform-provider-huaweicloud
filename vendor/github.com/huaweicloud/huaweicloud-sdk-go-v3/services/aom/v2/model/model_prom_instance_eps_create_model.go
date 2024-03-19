package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type PromInstanceEpsCreateModel struct {

	// Prometheus实例名称 名称不能以下划线或中划线开头结尾，只含有中文、英文、数字、下划线、中划线、长度1-100。
	PromName string `json:"prom_name"`

	// Prometheus实例id。
	PromId *string `json:"prom_id,omitempty"`

	// Prometheus实例类型（暂时不支持VPC、KUBERNETES）。
	PromType PromInstanceEpsCreateModelPromType `json:"prom_type"`

	// Prometheus实例版本号。
	PromVersion *string `json:"prom_version,omitempty"`

	// Prometheus实例创建时间戳。
	PromCreateTimestamp *int64 `json:"prom_create_timestamp,omitempty"`

	// Prometheus实例更新时间戳。
	PromUpdateTimestamp *int64 `json:"prom_update_timestamp,omitempty"`

	// Prometheus实例状态。
	PromStatus *PromInstanceEpsCreateModelPromStatus `json:"prom_status,omitempty"`

	// Prometheus实例所属的企业项目。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// Prometheus实例所属projectId。
	ProjectId *string `json:"project_id,omitempty"`

	// 删除时间。
	DeletedTime *int64 `json:"deleted_time,omitempty"`

	PromSpecConfig *PromConfigModel `json:"prom_spec_config,omitempty"`

	// Prometheus实例所属CCE特殊配置。
	CceSpecConfig *string `json:"cce_spec_config,omitempty"`
}

func (o PromInstanceEpsCreateModel) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PromInstanceEpsCreateModel struct{}"
	}

	return strings.Join([]string{"PromInstanceEpsCreateModel", string(data)}, " ")
}

type PromInstanceEpsCreateModelPromType struct {
	value string
}

type PromInstanceEpsCreateModelPromTypeEnum struct {
	ECS            PromInstanceEpsCreateModelPromType
	VPC            PromInstanceEpsCreateModelPromType
	CCE            PromInstanceEpsCreateModelPromType
	REMOTE_WRITE   PromInstanceEpsCreateModelPromType
	KUBERNETES     PromInstanceEpsCreateModelPromType
	CLOUD_SERVICE  PromInstanceEpsCreateModelPromType
	ACROSS_ACCOUNT PromInstanceEpsCreateModelPromType
}

func GetPromInstanceEpsCreateModelPromTypeEnum() PromInstanceEpsCreateModelPromTypeEnum {
	return PromInstanceEpsCreateModelPromTypeEnum{
		ECS: PromInstanceEpsCreateModelPromType{
			value: "ECS",
		},
		VPC: PromInstanceEpsCreateModelPromType{
			value: "VPC",
		},
		CCE: PromInstanceEpsCreateModelPromType{
			value: "CCE",
		},
		REMOTE_WRITE: PromInstanceEpsCreateModelPromType{
			value: "REMOTE_WRITE",
		},
		KUBERNETES: PromInstanceEpsCreateModelPromType{
			value: "KUBERNETES",
		},
		CLOUD_SERVICE: PromInstanceEpsCreateModelPromType{
			value: "CLOUD_SERVICE",
		},
		ACROSS_ACCOUNT: PromInstanceEpsCreateModelPromType{
			value: "ACROSS_ACCOUNT",
		},
	}
}

func (c PromInstanceEpsCreateModelPromType) Value() string {
	return c.value
}

func (c PromInstanceEpsCreateModelPromType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PromInstanceEpsCreateModelPromType) UnmarshalJSON(b []byte) error {
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

type PromInstanceEpsCreateModelPromStatus struct {
	value string
}

type PromInstanceEpsCreateModelPromStatusEnum struct {
	DELETED PromInstanceEpsCreateModelPromStatus
	NORMAL  PromInstanceEpsCreateModelPromStatus
	ALL     PromInstanceEpsCreateModelPromStatus
}

func GetPromInstanceEpsCreateModelPromStatusEnum() PromInstanceEpsCreateModelPromStatusEnum {
	return PromInstanceEpsCreateModelPromStatusEnum{
		DELETED: PromInstanceEpsCreateModelPromStatus{
			value: "DELETED",
		},
		NORMAL: PromInstanceEpsCreateModelPromStatus{
			value: "NORMAL",
		},
		ALL: PromInstanceEpsCreateModelPromStatus{
			value: "ALL",
		},
	}
}

func (c PromInstanceEpsCreateModelPromStatus) Value() string {
	return c.value
}

func (c PromInstanceEpsCreateModelPromStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PromInstanceEpsCreateModelPromStatus) UnmarshalJSON(b []byte) error {
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
