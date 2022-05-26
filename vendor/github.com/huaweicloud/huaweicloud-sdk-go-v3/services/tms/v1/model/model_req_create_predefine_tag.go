package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 创建预定义标签请求
type ReqCreatePredefineTag struct {

	// 操作标识（区分大小写）：create（创建）
	Action ReqCreatePredefineTagAction `json:"action"`

	// 标签列表
	Tags []PredefineTagRequest `json:"tags"`
}

func (o ReqCreatePredefineTag) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReqCreatePredefineTag struct{}"
	}

	return strings.Join([]string{"ReqCreatePredefineTag", string(data)}, " ")
}

type ReqCreatePredefineTagAction struct {
	value string
}

type ReqCreatePredefineTagActionEnum struct {
	CREATE ReqCreatePredefineTagAction
}

func GetReqCreatePredefineTagActionEnum() ReqCreatePredefineTagActionEnum {
	return ReqCreatePredefineTagActionEnum{
		CREATE: ReqCreatePredefineTagAction{
			value: "create",
		},
	}
}

func (c ReqCreatePredefineTagAction) Value() string {
	return c.value
}

func (c ReqCreatePredefineTagAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ReqCreatePredefineTagAction) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
