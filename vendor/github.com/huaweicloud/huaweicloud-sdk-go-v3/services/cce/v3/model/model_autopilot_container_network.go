package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// AutopilotContainerNetwork Container network parameters.
type AutopilotContainerNetwork struct {

	// 容器网络类型 - eni：云原生网络2.0，深度整合VPC原生ENI弹性网卡能力，采用VPC网段分配容器地址，支持ELB直通容器，享有高性能，创建集群时指定。
	Mode AutopilotContainerNetworkMode `json:"mode"`
}

func (o AutopilotContainerNetwork) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotContainerNetwork struct{}"
	}

	return strings.Join([]string{"AutopilotContainerNetwork", string(data)}, " ")
}

type AutopilotContainerNetworkMode struct {
	value string
}

type AutopilotContainerNetworkModeEnum struct {
	ENI AutopilotContainerNetworkMode
}

func GetAutopilotContainerNetworkModeEnum() AutopilotContainerNetworkModeEnum {
	return AutopilotContainerNetworkModeEnum{
		ENI: AutopilotContainerNetworkMode{
			value: "eni",
		},
	}
}

func (c AutopilotContainerNetworkMode) Value() string {
	return c.value
}

func (c AutopilotContainerNetworkMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AutopilotContainerNetworkMode) UnmarshalJSON(b []byte) error {
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
