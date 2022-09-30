package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListEngineProductsRequest struct {

	// 消息引擎的类型。
	Engine ListEngineProductsRequestEngine `json:"engine"`

	// 产品ID。
	ProductId *string `json:"product_id,omitempty"`
}

func (o ListEngineProductsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEngineProductsRequest struct{}"
	}

	return strings.Join([]string{"ListEngineProductsRequest", string(data)}, " ")
}

type ListEngineProductsRequestEngine struct {
	value string
}

type ListEngineProductsRequestEngineEnum struct {
	KAFKA ListEngineProductsRequestEngine
}

func GetListEngineProductsRequestEngineEnum() ListEngineProductsRequestEngineEnum {
	return ListEngineProductsRequestEngineEnum{
		KAFKA: ListEngineProductsRequestEngine{
			value: "kafka",
		},
	}
}

func (c ListEngineProductsRequestEngine) Value() string {
	return c.value
}

func (c ListEngineProductsRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListEngineProductsRequestEngine) UnmarshalJSON(b []byte) error {
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
