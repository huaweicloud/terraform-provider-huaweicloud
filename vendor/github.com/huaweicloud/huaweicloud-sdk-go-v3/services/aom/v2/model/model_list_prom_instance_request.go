package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListPromInstanceRequest Request Object
type ListPromInstanceRequest struct {

	// Prometheus实例id(prom_id与prom_type同时存在时，仅prom_id生效)。
	PromId *string `json:"prom_id,omitempty"`

	// Prometheus实例类型（暂时不支持VPC、KUBERNETES）。
	PromType *ListPromInstanceRequestPromType `json:"prom_type,omitempty"`

	// cce集群开关。
	CceClusterEnable *ListPromInstanceRequestCceClusterEnable `json:"cce_cluster_enable,omitempty"`

	// Prometheus实例状态。
	PromStatus *ListPromInstanceRequestPromStatus `json:"prom_status,omitempty"`

	// 企业项目id。 - 查询单个企业项目下实例，填写企业项目id。 - 查询所有企业项目下实例，填写“all_granted_eps”。
	EnterpriseProjectId *string `json:"Enterprise-Project-Id,omitempty"`
}

func (o ListPromInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPromInstanceRequest struct{}"
	}

	return strings.Join([]string{"ListPromInstanceRequest", string(data)}, " ")
}

type ListPromInstanceRequestPromType struct {
	value string
}

type ListPromInstanceRequestPromTypeEnum struct {
	DEFAULT        ListPromInstanceRequestPromType
	ECS            ListPromInstanceRequestPromType
	VPC            ListPromInstanceRequestPromType
	CCE            ListPromInstanceRequestPromType
	REMOTE_WRITE   ListPromInstanceRequestPromType
	KUBERNETES     ListPromInstanceRequestPromType
	CLOUD_SERVICE  ListPromInstanceRequestPromType
	ACROSS_ACCOUNT ListPromInstanceRequestPromType
}

func GetListPromInstanceRequestPromTypeEnum() ListPromInstanceRequestPromTypeEnum {
	return ListPromInstanceRequestPromTypeEnum{
		DEFAULT: ListPromInstanceRequestPromType{
			value: "default",
		},
		ECS: ListPromInstanceRequestPromType{
			value: "ECS",
		},
		VPC: ListPromInstanceRequestPromType{
			value: "VPC",
		},
		CCE: ListPromInstanceRequestPromType{
			value: "CCE",
		},
		REMOTE_WRITE: ListPromInstanceRequestPromType{
			value: "REMOTE_WRITE",
		},
		KUBERNETES: ListPromInstanceRequestPromType{
			value: "KUBERNETES",
		},
		CLOUD_SERVICE: ListPromInstanceRequestPromType{
			value: "CLOUD_SERVICE",
		},
		ACROSS_ACCOUNT: ListPromInstanceRequestPromType{
			value: "ACROSS_ACCOUNT",
		},
	}
}

func (c ListPromInstanceRequestPromType) Value() string {
	return c.value
}

func (c ListPromInstanceRequestPromType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListPromInstanceRequestPromType) UnmarshalJSON(b []byte) error {
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

type ListPromInstanceRequestCceClusterEnable struct {
	value string
}

type ListPromInstanceRequestCceClusterEnableEnum struct {
	TRUE  ListPromInstanceRequestCceClusterEnable
	FALSE ListPromInstanceRequestCceClusterEnable
}

func GetListPromInstanceRequestCceClusterEnableEnum() ListPromInstanceRequestCceClusterEnableEnum {
	return ListPromInstanceRequestCceClusterEnableEnum{
		TRUE: ListPromInstanceRequestCceClusterEnable{
			value: "true",
		},
		FALSE: ListPromInstanceRequestCceClusterEnable{
			value: "false",
		},
	}
}

func (c ListPromInstanceRequestCceClusterEnable) Value() string {
	return c.value
}

func (c ListPromInstanceRequestCceClusterEnable) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListPromInstanceRequestCceClusterEnable) UnmarshalJSON(b []byte) error {
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

type ListPromInstanceRequestPromStatus struct {
	value string
}

type ListPromInstanceRequestPromStatusEnum struct {
	DELETED ListPromInstanceRequestPromStatus
	NORMAL  ListPromInstanceRequestPromStatus
	ALL     ListPromInstanceRequestPromStatus
}

func GetListPromInstanceRequestPromStatusEnum() ListPromInstanceRequestPromStatusEnum {
	return ListPromInstanceRequestPromStatusEnum{
		DELETED: ListPromInstanceRequestPromStatus{
			value: "DELETED",
		},
		NORMAL: ListPromInstanceRequestPromStatus{
			value: "NORMAL",
		},
		ALL: ListPromInstanceRequestPromStatus{
			value: "ALL",
		},
	}
}

func (c ListPromInstanceRequestPromStatus) Value() string {
	return c.value
}

func (c ListPromInstanceRequestPromStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListPromInstanceRequestPromStatus) UnmarshalJSON(b []byte) error {
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
