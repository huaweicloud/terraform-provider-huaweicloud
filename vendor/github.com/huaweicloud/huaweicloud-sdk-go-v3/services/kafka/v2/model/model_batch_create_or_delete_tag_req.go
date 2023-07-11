package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type BatchCreateOrDeleteTagReq struct {

	// 操作标识（仅支持小写）: - create（创建） - delete（删除）
	Action *BatchCreateOrDeleteTagReqAction `json:"action,omitempty"`

	// 标签列表。
	Tags *[]TagEntity `json:"tags,omitempty"`
}

func (o BatchCreateOrDeleteTagReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateOrDeleteTagReq struct{}"
	}

	return strings.Join([]string{"BatchCreateOrDeleteTagReq", string(data)}, " ")
}

type BatchCreateOrDeleteTagReqAction struct {
	value string
}

type BatchCreateOrDeleteTagReqActionEnum struct {
	CREATE BatchCreateOrDeleteTagReqAction
	DELETE BatchCreateOrDeleteTagReqAction
}

func GetBatchCreateOrDeleteTagReqActionEnum() BatchCreateOrDeleteTagReqActionEnum {
	return BatchCreateOrDeleteTagReqActionEnum{
		CREATE: BatchCreateOrDeleteTagReqAction{
			value: "create",
		},
		DELETE: BatchCreateOrDeleteTagReqAction{
			value: "delete",
		},
	}
}

func (c BatchCreateOrDeleteTagReqAction) Value() string {
	return c.value
}

func (c BatchCreateOrDeleteTagReqAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BatchCreateOrDeleteTagReqAction) UnmarshalJSON(b []byte) error {
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
