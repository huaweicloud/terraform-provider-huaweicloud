package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListInstancesRecommendationRequest Request Object
type ListInstancesRecommendationRequest struct {

	// 引擎类型
	Engine ListInstancesRecommendationRequestEngine `json:"engine"`
}

func (o ListInstancesRecommendationRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstancesRecommendationRequest struct{}"
	}

	return strings.Join([]string{"ListInstancesRecommendationRequest", string(data)}, " ")
}

type ListInstancesRecommendationRequestEngine struct {
	value string
}

type ListInstancesRecommendationRequestEngineEnum struct {
	MYSQL      ListInstancesRecommendationRequestEngine
	POSTGRESQL ListInstancesRecommendationRequestEngine
	SQLSERVER  ListInstancesRecommendationRequestEngine
}

func GetListInstancesRecommendationRequestEngineEnum() ListInstancesRecommendationRequestEngineEnum {
	return ListInstancesRecommendationRequestEngineEnum{
		MYSQL: ListInstancesRecommendationRequestEngine{
			value: "mysql",
		},
		POSTGRESQL: ListInstancesRecommendationRequestEngine{
			value: "postgresql",
		},
		SQLSERVER: ListInstancesRecommendationRequestEngine{
			value: "sqlserver",
		},
	}
}

func (c ListInstancesRecommendationRequestEngine) Value() string {
	return c.value
}

func (c ListInstancesRecommendationRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstancesRecommendationRequestEngine) UnmarshalJSON(b []byte) error {
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
