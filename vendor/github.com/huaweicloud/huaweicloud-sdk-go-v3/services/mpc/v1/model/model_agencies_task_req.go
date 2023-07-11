package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type AgenciesTaskReq struct {

	// 委托任务租户Id
	ProjectId *string `json:"project_id,omitempty"`

	// 操作标记，取值[CREATED,CANCELED]，CREATED表示授权, CANCELED表示取消授权
	OperateType AgenciesTaskReqOperateType `json:"operate_type"`
}

func (o AgenciesTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AgenciesTaskReq struct{}"
	}

	return strings.Join([]string{"AgenciesTaskReq", string(data)}, " ")
}

type AgenciesTaskReqOperateType struct {
	value string
}

type AgenciesTaskReqOperateTypeEnum struct {
	CREATED  AgenciesTaskReqOperateType
	CANCELED AgenciesTaskReqOperateType
}

func GetAgenciesTaskReqOperateTypeEnum() AgenciesTaskReqOperateTypeEnum {
	return AgenciesTaskReqOperateTypeEnum{
		CREATED: AgenciesTaskReqOperateType{
			value: "CREATED",
		},
		CANCELED: AgenciesTaskReqOperateType{
			value: "CANCELED",
		},
	}
}

func (c AgenciesTaskReqOperateType) Value() string {
	return c.value
}

func (c AgenciesTaskReqOperateType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AgenciesTaskReqOperateType) UnmarshalJSON(b []byte) error {
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
