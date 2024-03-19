package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type BatchDeleteInstanceUsersReq struct {

	// 删除类型。当前只支持delete。  当删除用户时，为必选参数。
	Action *BatchDeleteInstanceUsersReqAction `json:"action,omitempty"`

	// 用户列表。  当删除用户时，为必选参数。
	Users *[]string `json:"users,omitempty"`
}

func (o BatchDeleteInstanceUsersReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteInstanceUsersReq struct{}"
	}

	return strings.Join([]string{"BatchDeleteInstanceUsersReq", string(data)}, " ")
}

type BatchDeleteInstanceUsersReqAction struct {
	value string
}

type BatchDeleteInstanceUsersReqActionEnum struct {
	DELETE BatchDeleteInstanceUsersReqAction
}

func GetBatchDeleteInstanceUsersReqActionEnum() BatchDeleteInstanceUsersReqActionEnum {
	return BatchDeleteInstanceUsersReqActionEnum{
		DELETE: BatchDeleteInstanceUsersReqAction{
			value: "delete",
		},
	}
}

func (c BatchDeleteInstanceUsersReqAction) Value() string {
	return c.value
}

func (c BatchDeleteInstanceUsersReqAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BatchDeleteInstanceUsersReqAction) UnmarshalJSON(b []byte) error {
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
