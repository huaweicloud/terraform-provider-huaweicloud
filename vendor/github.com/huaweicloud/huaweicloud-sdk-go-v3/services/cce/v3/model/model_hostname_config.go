package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type HostnameConfig struct {

	// K8S节点名称配置类型, 默认为“privateIp”。  -  privateIp: 将节点私有IP作为K8S节点名称 -  cceNodeName: 将CCE节点名称作为K8S节点名称  > - 配置为cceNodeName的节点, 其节点名称、K8S节点名称以及虚机名称相同。节点名称不支持修改, 并且在ECS侧修改了虚机名称，同步云服务器时，不会将修改后的虚机名称同步到节点。 > - 配置为cceNodeName的节点，为了避免K8S节点名称冲突，系统会自动在节点名称后添加后缀，后缀的格式为中划线(-)+五位随机字符，随机字符的取值为[a-z0-9]。
	Type HostnameConfigType `json:"type"`
}

func (o HostnameConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HostnameConfig struct{}"
	}

	return strings.Join([]string{"HostnameConfig", string(data)}, " ")
}

type HostnameConfigType struct {
	value string
}

type HostnameConfigTypeEnum struct {
	PRIVATE_IP    HostnameConfigType
	CCE_NODE_NAME HostnameConfigType
}

func GetHostnameConfigTypeEnum() HostnameConfigTypeEnum {
	return HostnameConfigTypeEnum{
		PRIVATE_IP: HostnameConfigType{
			value: "privateIp",
		},
		CCE_NODE_NAME: HostnameConfigType{
			value: "cceNodeName",
		},
	}
}

func (c HostnameConfigType) Value() string {
	return c.value
}

func (c HostnameConfigType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *HostnameConfigType) UnmarshalJSON(b []byte) error {
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
