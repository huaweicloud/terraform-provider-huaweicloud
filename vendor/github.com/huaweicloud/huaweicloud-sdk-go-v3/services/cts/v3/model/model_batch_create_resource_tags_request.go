package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// BatchCreateResourceTagsRequest Request Object
type BatchCreateResourceTagsRequest struct {

	// 资源ID。
	ResourceId string `json:"resource_id"`

	// CTS服务的资源类型，目前仅支持cts-tracker。
	ResourceType BatchCreateResourceTagsRequestResourceType `json:"resource_type"`

	Body *BatchCreateResourceTagsRequestBody `json:"body,omitempty"`
}

func (o BatchCreateResourceTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateResourceTagsRequest struct{}"
	}

	return strings.Join([]string{"BatchCreateResourceTagsRequest", string(data)}, " ")
}

type BatchCreateResourceTagsRequestResourceType struct {
	value string
}

type BatchCreateResourceTagsRequestResourceTypeEnum struct {
	CTS_TRACKER BatchCreateResourceTagsRequestResourceType
}

func GetBatchCreateResourceTagsRequestResourceTypeEnum() BatchCreateResourceTagsRequestResourceTypeEnum {
	return BatchCreateResourceTagsRequestResourceTypeEnum{
		CTS_TRACKER: BatchCreateResourceTagsRequestResourceType{
			value: "cts-tracker",
		},
	}
}

func (c BatchCreateResourceTagsRequestResourceType) Value() string {
	return c.value
}

func (c BatchCreateResourceTagsRequestResourceType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BatchCreateResourceTagsRequestResourceType) UnmarshalJSON(b []byte) error {
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
