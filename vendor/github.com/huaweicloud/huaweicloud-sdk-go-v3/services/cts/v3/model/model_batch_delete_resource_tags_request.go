package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// BatchDeleteResourceTagsRequest Request Object
type BatchDeleteResourceTagsRequest struct {

	// 资源ID。
	ResourceId string `json:"resource_id"`

	// CTS服务的资源类型，目前仅支持cts-tracker。
	ResourceType BatchDeleteResourceTagsRequestResourceType `json:"resource_type"`

	Body *BatchDeleteResourceTagsRequestBody `json:"body,omitempty"`
}

func (o BatchDeleteResourceTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteResourceTagsRequest struct{}"
	}

	return strings.Join([]string{"BatchDeleteResourceTagsRequest", string(data)}, " ")
}

type BatchDeleteResourceTagsRequestResourceType struct {
	value string
}

type BatchDeleteResourceTagsRequestResourceTypeEnum struct {
	CTS_TRACKER BatchDeleteResourceTagsRequestResourceType
}

func GetBatchDeleteResourceTagsRequestResourceTypeEnum() BatchDeleteResourceTagsRequestResourceTypeEnum {
	return BatchDeleteResourceTagsRequestResourceTypeEnum{
		CTS_TRACKER: BatchDeleteResourceTagsRequestResourceType{
			value: "cts-tracker",
		},
	}
}

func (c BatchDeleteResourceTagsRequestResourceType) Value() string {
	return c.value
}

func (c BatchDeleteResourceTagsRequestResourceType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BatchDeleteResourceTagsRequestResourceType) UnmarshalJSON(b []byte) error {
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
