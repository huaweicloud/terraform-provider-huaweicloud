package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// UpdateAzByInstanceTypeRequest Request Object
type UpdateAzByInstanceTypeRequest struct {

	// 待切换AZ的集群ID。
	ClusterId string `json:"cluster_id"`

	// 待切换AZ的节点类型。支持: - all：所有节点类型。 - ess： 数据节点。 - ess-cold: 冷数据节点。 - ess-client: Client节点。 - ess-master: Master节点。
	InstType UpdateAzByInstanceTypeRequestInstType `json:"inst_type"`

	Body *UpdateAzByInstanceTypeReq `json:"body,omitempty"`
}

func (o UpdateAzByInstanceTypeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAzByInstanceTypeRequest struct{}"
	}

	return strings.Join([]string{"UpdateAzByInstanceTypeRequest", string(data)}, " ")
}

type UpdateAzByInstanceTypeRequestInstType struct {
	value string
}

type UpdateAzByInstanceTypeRequestInstTypeEnum struct {
	ALL        UpdateAzByInstanceTypeRequestInstType
	ESS        UpdateAzByInstanceTypeRequestInstType
	ESS_COLD   UpdateAzByInstanceTypeRequestInstType
	ESS_CLIENT UpdateAzByInstanceTypeRequestInstType
	ESS_MASTER UpdateAzByInstanceTypeRequestInstType
}

func GetUpdateAzByInstanceTypeRequestInstTypeEnum() UpdateAzByInstanceTypeRequestInstTypeEnum {
	return UpdateAzByInstanceTypeRequestInstTypeEnum{
		ALL: UpdateAzByInstanceTypeRequestInstType{
			value: "all",
		},
		ESS: UpdateAzByInstanceTypeRequestInstType{
			value: "ess",
		},
		ESS_COLD: UpdateAzByInstanceTypeRequestInstType{
			value: "ess-cold",
		},
		ESS_CLIENT: UpdateAzByInstanceTypeRequestInstType{
			value: "ess-client",
		},
		ESS_MASTER: UpdateAzByInstanceTypeRequestInstType{
			value: "ess-master",
		},
	}
}

func (c UpdateAzByInstanceTypeRequestInstType) Value() string {
	return c.value
}

func (c UpdateAzByInstanceTypeRequestInstType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateAzByInstanceTypeRequestInstType) UnmarshalJSON(b []byte) error {
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
