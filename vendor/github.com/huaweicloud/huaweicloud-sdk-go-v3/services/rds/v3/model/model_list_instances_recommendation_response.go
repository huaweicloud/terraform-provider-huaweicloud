package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListInstancesRecommendationResponse Response Object
type ListInstancesRecommendationResponse struct {

	// 引擎类型
	Engine *ListInstancesRecommendationResponseEngine `json:"engine,omitempty"`

	// 推荐商品信息
	RecommendedProducts *[]RecommendedProduct `json:"recommended_products,omitempty"`
	HttpStatusCode      int                   `json:"-"`
}

func (o ListInstancesRecommendationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstancesRecommendationResponse struct{}"
	}

	return strings.Join([]string{"ListInstancesRecommendationResponse", string(data)}, " ")
}

type ListInstancesRecommendationResponseEngine struct {
	value string
}

type ListInstancesRecommendationResponseEngineEnum struct {
	MYSQL      ListInstancesRecommendationResponseEngine
	POSTGRESQL ListInstancesRecommendationResponseEngine
	SQLSERVER  ListInstancesRecommendationResponseEngine
}

func GetListInstancesRecommendationResponseEngineEnum() ListInstancesRecommendationResponseEngineEnum {
	return ListInstancesRecommendationResponseEngineEnum{
		MYSQL: ListInstancesRecommendationResponseEngine{
			value: "mysql",
		},
		POSTGRESQL: ListInstancesRecommendationResponseEngine{
			value: "postgresql",
		},
		SQLSERVER: ListInstancesRecommendationResponseEngine{
			value: "sqlserver",
		},
	}
}

func (c ListInstancesRecommendationResponseEngine) Value() string {
	return c.value
}

func (c ListInstancesRecommendationResponseEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstancesRecommendationResponseEngine) UnmarshalJSON(b []byte) error {
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
