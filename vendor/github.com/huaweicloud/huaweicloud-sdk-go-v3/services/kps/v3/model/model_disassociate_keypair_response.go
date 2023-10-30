package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// DisassociateKeypairResponse Response Object
type DisassociateKeypairResponse struct {

	// 任务下发成功返回的ID。
	TaskId *string `json:"task_id,omitempty"`

	// 绑定的虚拟机id。
	ServerId *string `json:"server_id,omitempty"`

	// 任务下发的状态。SUCCESS或FAILED。
	Status *DisassociateKeypairResponseStatus `json:"status,omitempty"`

	// 任务下发失败返回的错误码。
	ErrorCode *string `json:"error_code,omitempty"`

	// 任务下发失败返回的错误信息。
	ErrorMsg       *string `json:"error_msg,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DisassociateKeypairResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DisassociateKeypairResponse struct{}"
	}

	return strings.Join([]string{"DisassociateKeypairResponse", string(data)}, " ")
}

type DisassociateKeypairResponseStatus struct {
	value string
}

type DisassociateKeypairResponseStatusEnum struct {
	SUCCESS DisassociateKeypairResponseStatus
	FAILED  DisassociateKeypairResponseStatus
}

func GetDisassociateKeypairResponseStatusEnum() DisassociateKeypairResponseStatusEnum {
	return DisassociateKeypairResponseStatusEnum{
		SUCCESS: DisassociateKeypairResponseStatus{
			value: "SUCCESS",
		},
		FAILED: DisassociateKeypairResponseStatus{
			value: "FAILED",
		},
	}
}

func (c DisassociateKeypairResponseStatus) Value() string {
	return c.value
}

func (c DisassociateKeypairResponseStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DisassociateKeypairResponseStatus) UnmarshalJSON(b []byte) error {
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
