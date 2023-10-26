package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListProductsRequest Request Object
type ListProductsRequest struct {

	// 消息引擎的类型。当前只支持kafka类型。
	Engine ListProductsRequestEngine `json:"engine"`
}

func (o ListProductsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProductsRequest struct{}"
	}

	return strings.Join([]string{"ListProductsRequest", string(data)}, " ")
}

type ListProductsRequestEngine struct {
	value string
}

type ListProductsRequestEngineEnum struct {
	KAFKA ListProductsRequestEngine
}

func GetListProductsRequestEngineEnum() ListProductsRequestEngineEnum {
	return ListProductsRequestEngineEnum{
		KAFKA: ListProductsRequestEngine{
			value: "kafka",
		},
	}
}

func (c ListProductsRequestEngine) Value() string {
	return c.value
}

func (c ListProductsRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListProductsRequestEngine) UnmarshalJSON(b []byte) error {
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
