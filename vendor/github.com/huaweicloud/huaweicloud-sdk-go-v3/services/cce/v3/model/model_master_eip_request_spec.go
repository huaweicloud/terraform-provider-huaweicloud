package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type MasterEipRequestSpec struct {

	// 绑定或解绑动作，必选参数。 - 绑定：固定值为{\"action\":\"bind\"} - 解绑：固定值为{\"action\":\"unbind\"}
	Action *MasterEipRequestSpecAction `json:"action,omitempty"`

	Spec *MasterEipRequestSpecSpec `json:"spec,omitempty"`

	// 带宽(字段已失效，暂不推荐使用)
	Bandwidth *string `json:"bandwidth,omitempty"`

	// 弹性网卡IP(字段已失效，暂不推荐使用)
	ElasticIp *string `json:"elasticIp,omitempty"`
}

func (o MasterEipRequestSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MasterEipRequestSpec struct{}"
	}

	return strings.Join([]string{"MasterEipRequestSpec", string(data)}, " ")
}

type MasterEipRequestSpecAction struct {
	value string
}

type MasterEipRequestSpecActionEnum struct {
	BIND   MasterEipRequestSpecAction
	UNBIND MasterEipRequestSpecAction
}

func GetMasterEipRequestSpecActionEnum() MasterEipRequestSpecActionEnum {
	return MasterEipRequestSpecActionEnum{
		BIND: MasterEipRequestSpecAction{
			value: "bind",
		},
		UNBIND: MasterEipRequestSpecAction{
			value: "unbind",
		},
	}
}

func (c MasterEipRequestSpecAction) Value() string {
	return c.value
}

func (c MasterEipRequestSpecAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MasterEipRequestSpecAction) UnmarshalJSON(b []byte) error {
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
