package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 专属资源池信息。
type DedicatedResource struct {
	// 专属资源池ID。

	Id *string `json:"id,omitempty"`
	// 专属资源池名称

	ResourceName *string `json:"resource_name,omitempty"`
	// 数据库引擎名称

	EngineName *string `json:"engine_name,omitempty"`
	// CPU架构

	Architecture *string `json:"architecture,omitempty"`
	// 专属资源池状态

	Status *DedicatedResourceStatus `json:"status,omitempty"`

	Capacity *DedicatedResourceCapacity `json:"capacity,omitempty"`
	// 专属资源池可用区信息。

	AvailabilityZone *[]string `json:"availability_zone,omitempty"`
}

func (o DedicatedResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DedicatedResource struct{}"
	}

	return strings.Join([]string{"DedicatedResource", string(data)}, " ")
}

type DedicatedResourceStatus struct {
	value string
}

type DedicatedResourceStatusEnum struct {
	NORMAL    DedicatedResourceStatus
	BUILDING  DedicatedResourceStatus
	EXTENDING DedicatedResourceStatus
	DELETED   DedicatedResourceStatus
}

func GetDedicatedResourceStatusEnum() DedicatedResourceStatusEnum {
	return DedicatedResourceStatusEnum{
		NORMAL: DedicatedResourceStatus{
			value: "NORMAL",
		},
		BUILDING: DedicatedResourceStatus{
			value: "BUILDING",
		},
		EXTENDING: DedicatedResourceStatus{
			value: "EXTENDING",
		},
		DELETED: DedicatedResourceStatus{
			value: "DELETED",
		},
	}
}

func (c DedicatedResourceStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DedicatedResourceStatus) UnmarshalJSON(b []byte) error {
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
