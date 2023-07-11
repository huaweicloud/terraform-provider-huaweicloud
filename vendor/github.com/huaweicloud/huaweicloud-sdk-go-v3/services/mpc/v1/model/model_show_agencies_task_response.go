package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowAgenciesTaskResponse Response Object
type ShowAgenciesTaskResponse struct {

	// 操作标记，取值[CREATED,CANCELED]，CREATED表示授权, CANCELED表示取消授权
	OperateType    *ShowAgenciesTaskResponseOperateType `json:"operate_type,omitempty"`
	HttpStatusCode int                                  `json:"-"`
}

func (o ShowAgenciesTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAgenciesTaskResponse struct{}"
	}

	return strings.Join([]string{"ShowAgenciesTaskResponse", string(data)}, " ")
}

type ShowAgenciesTaskResponseOperateType struct {
	value string
}

type ShowAgenciesTaskResponseOperateTypeEnum struct {
	CREATED  ShowAgenciesTaskResponseOperateType
	CANCELED ShowAgenciesTaskResponseOperateType
}

func GetShowAgenciesTaskResponseOperateTypeEnum() ShowAgenciesTaskResponseOperateTypeEnum {
	return ShowAgenciesTaskResponseOperateTypeEnum{
		CREATED: ShowAgenciesTaskResponseOperateType{
			value: "CREATED",
		},
		CANCELED: ShowAgenciesTaskResponseOperateType{
			value: "CANCELED",
		},
	}
}

func (c ShowAgenciesTaskResponseOperateType) Value() string {
	return c.value
}

func (c ShowAgenciesTaskResponseOperateType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowAgenciesTaskResponseOperateType) UnmarshalJSON(b []byte) error {
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
