package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowClusterConfigRequest Request Object
type ShowClusterConfigRequest struct {

	// 组件类型 , 合法取值为control，audit，system-addon。不填写则查询全部类型。 - control 控制面组件日志。 - audit 控制面审计日志。 - system-addon 系统插件日志。
	Type *ShowClusterConfigRequestType `json:"type,omitempty"`

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`
}

func (o ShowClusterConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterConfigRequest struct{}"
	}

	return strings.Join([]string{"ShowClusterConfigRequest", string(data)}, " ")
}

type ShowClusterConfigRequestType struct {
	value string
}

type ShowClusterConfigRequestTypeEnum struct {
	CONTROL      ShowClusterConfigRequestType
	AUDIT        ShowClusterConfigRequestType
	SYSTEM_ADDON ShowClusterConfigRequestType
}

func GetShowClusterConfigRequestTypeEnum() ShowClusterConfigRequestTypeEnum {
	return ShowClusterConfigRequestTypeEnum{
		CONTROL: ShowClusterConfigRequestType{
			value: "control",
		},
		AUDIT: ShowClusterConfigRequestType{
			value: "audit",
		},
		SYSTEM_ADDON: ShowClusterConfigRequestType{
			value: "system-addon",
		},
	}
}

func (c ShowClusterConfigRequestType) Value() string {
	return c.value
}

func (c ShowClusterConfigRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowClusterConfigRequestType) UnmarshalJSON(b []byte) error {
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
