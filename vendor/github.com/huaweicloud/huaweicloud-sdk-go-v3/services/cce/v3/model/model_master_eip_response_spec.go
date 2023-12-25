package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type MasterEipResponseSpec struct {

	// 绑定动作
	Action *MasterEipResponseSpecAction `json:"action,omitempty"`

	Spec *MasterEipResponseSpecSpec `json:"spec,omitempty"`

	// 弹性公网IP
	ElasticIp *string `json:"elasticIp,omitempty"`
}

func (o MasterEipResponseSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MasterEipResponseSpec struct{}"
	}

	return strings.Join([]string{"MasterEipResponseSpec", string(data)}, " ")
}

type MasterEipResponseSpecAction struct {
	value string
}

type MasterEipResponseSpecActionEnum struct {
	BIND MasterEipResponseSpecAction
}

func GetMasterEipResponseSpecActionEnum() MasterEipResponseSpecActionEnum {
	return MasterEipResponseSpecActionEnum{
		BIND: MasterEipResponseSpecAction{
			value: "bind",
		},
	}
}

func (c MasterEipResponseSpecAction) Value() string {
	return c.value
}

func (c MasterEipResponseSpecAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MasterEipResponseSpecAction) UnmarshalJSON(b []byte) error {
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
