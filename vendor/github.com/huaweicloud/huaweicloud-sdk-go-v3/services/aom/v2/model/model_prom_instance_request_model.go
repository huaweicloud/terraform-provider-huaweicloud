package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type PromInstanceRequestModel struct {

	// Prometheus实例名称 名称不能以下划线或中划线开头结尾，只含有中文、英文、数字、下划线、中划线、长度1-100。
	PromName string `json:"prom_name"`

	// Prometheus实例类型（暂时不支持VPC、KUBERNETES）。
	PromType PromInstanceRequestModelPromType `json:"prom_type"`

	// Prometheus实例版本号。
	PromVersion *string `json:"prom_version,omitempty"`

	// Prometheus实例所属的企业项目。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// Prometheus实例所属projectId。
	ProjectId *string `json:"project_id,omitempty"`
}

func (o PromInstanceRequestModel) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PromInstanceRequestModel struct{}"
	}

	return strings.Join([]string{"PromInstanceRequestModel", string(data)}, " ")
}

type PromInstanceRequestModelPromType struct {
	value string
}

type PromInstanceRequestModelPromTypeEnum struct {
	ECS            PromInstanceRequestModelPromType
	VPC            PromInstanceRequestModelPromType
	CCE            PromInstanceRequestModelPromType
	REMOTE_WRITE   PromInstanceRequestModelPromType
	KUBERNETES     PromInstanceRequestModelPromType
	CLOUD_SERVICE  PromInstanceRequestModelPromType
	ACROSS_ACCOUNT PromInstanceRequestModelPromType
}

func GetPromInstanceRequestModelPromTypeEnum() PromInstanceRequestModelPromTypeEnum {
	return PromInstanceRequestModelPromTypeEnum{
		ECS: PromInstanceRequestModelPromType{
			value: "ECS",
		},
		VPC: PromInstanceRequestModelPromType{
			value: "VPC",
		},
		CCE: PromInstanceRequestModelPromType{
			value: "CCE",
		},
		REMOTE_WRITE: PromInstanceRequestModelPromType{
			value: "REMOTE_WRITE",
		},
		KUBERNETES: PromInstanceRequestModelPromType{
			value: "KUBERNETES",
		},
		CLOUD_SERVICE: PromInstanceRequestModelPromType{
			value: "CLOUD_SERVICE",
		},
		ACROSS_ACCOUNT: PromInstanceRequestModelPromType{
			value: "ACROSS_ACCOUNT",
		},
	}
}

func (c PromInstanceRequestModelPromType) Value() string {
	return c.value
}

func (c PromInstanceRequestModelPromType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PromInstanceRequestModelPromType) UnmarshalJSON(b []byte) error {
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
