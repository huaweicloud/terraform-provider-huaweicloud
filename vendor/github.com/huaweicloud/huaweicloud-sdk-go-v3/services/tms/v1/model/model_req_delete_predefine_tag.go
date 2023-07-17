package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ReqDeletePredefineTag 删除预定义标签请求
type ReqDeletePredefineTag struct {

	// 操作标识（区分大小写）：delete（删除）
	Action ReqDeletePredefineTagAction `json:"action"`

	// 标签列表
	Tags []PredefineTagRequest `json:"tags"`
}

func (o ReqDeletePredefineTag) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReqDeletePredefineTag struct{}"
	}

	return strings.Join([]string{"ReqDeletePredefineTag", string(data)}, " ")
}

type ReqDeletePredefineTagAction struct {
	value string
}

type ReqDeletePredefineTagActionEnum struct {
	DELETE ReqDeletePredefineTagAction
}

func GetReqDeletePredefineTagActionEnum() ReqDeletePredefineTagActionEnum {
	return ReqDeletePredefineTagActionEnum{
		DELETE: ReqDeletePredefineTagAction{
			value: "delete",
		},
	}
}

func (c ReqDeletePredefineTagAction) Value() string {
	return c.value
}

func (c ReqDeletePredefineTagAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ReqDeletePredefineTagAction) UnmarshalJSON(b []byte) error {
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
