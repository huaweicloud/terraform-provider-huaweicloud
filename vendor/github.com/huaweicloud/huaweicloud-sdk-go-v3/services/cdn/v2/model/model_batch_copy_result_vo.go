package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// BatchCopyResultVo 成功响应详细内容。
type BatchCopyResultVo struct {

	// 失败原因,成功时没有该字段
	Reason *string `json:"reason,omitempty"`

	// 批量操作结果。
	Status BatchCopyResultVoStatus `json:"status"`

	// 域名。
	DomainName string `json:"domain_name"`
}

func (o BatchCopyResultVo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCopyResultVo struct{}"
	}

	return strings.Join([]string{"BatchCopyResultVo", string(data)}, " ")
}

type BatchCopyResultVoStatus struct {
	value string
}

type BatchCopyResultVoStatusEnum struct {
	SUCCESS BatchCopyResultVoStatus
	FAIL    BatchCopyResultVoStatus
}

func GetBatchCopyResultVoStatusEnum() BatchCopyResultVoStatusEnum {
	return BatchCopyResultVoStatusEnum{
		SUCCESS: BatchCopyResultVoStatus{
			value: "success",
		},
		FAIL: BatchCopyResultVoStatus{
			value: "fail",
		},
	}
}

func (c BatchCopyResultVoStatus) Value() string {
	return c.value
}

func (c BatchCopyResultVoStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BatchCopyResultVoStatus) UnmarshalJSON(b []byte) error {
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
