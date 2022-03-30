package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type Resource struct {
	// 指定类型的配额。 - instance: 表示实例的配额

	Type ResourceType `json:"type"`
	// 已创建的资源个数。

	Used int32 `json:"used"`
	// 资源最大的配额数。

	Quota int32 `json:"quota"`
}

func (o Resource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Resource struct{}"
	}

	return strings.Join([]string{"Resource", string(data)}, " ")
}

type ResourceType struct {
	value string
}

type ResourceTypeEnum struct {
	INSTANCE ResourceType
}

func GetResourceTypeEnum() ResourceTypeEnum {
	return ResourceTypeEnum{
		INSTANCE: ResourceType{
			value: "instance",
		},
	}
}

func (c ResourceType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ResourceType) UnmarshalJSON(b []byte) error {
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
