package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// UpdateReleaseReqBody 更新模板实例的请求体
type UpdateReleaseReqBody struct {

	// 模板ID
	ChartId string `json:"chart_id"`

	// 更新操作，升级为upgrade，回退为rollback
	Action UpdateReleaseReqBodyAction `json:"action"`

	Parameters *ReleaseReqBodyParams `json:"parameters"`

	Values *CreateReleaseReqBodyValues `json:"values"`
}

func (o UpdateReleaseReqBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateReleaseReqBody struct{}"
	}

	return strings.Join([]string{"UpdateReleaseReqBody", string(data)}, " ")
}

type UpdateReleaseReqBodyAction struct {
	value string
}

type UpdateReleaseReqBodyActionEnum struct {
	UPGRADE  UpdateReleaseReqBodyAction
	ROLLBACK UpdateReleaseReqBodyAction
}

func GetUpdateReleaseReqBodyActionEnum() UpdateReleaseReqBodyActionEnum {
	return UpdateReleaseReqBodyActionEnum{
		UPGRADE: UpdateReleaseReqBodyAction{
			value: "upgrade",
		},
		ROLLBACK: UpdateReleaseReqBodyAction{
			value: "rollback",
		},
	}
}

func (c UpdateReleaseReqBodyAction) Value() string {
	return c.value
}

func (c UpdateReleaseReqBodyAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateReleaseReqBodyAction) UnmarshalJSON(b []byte) error {
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
