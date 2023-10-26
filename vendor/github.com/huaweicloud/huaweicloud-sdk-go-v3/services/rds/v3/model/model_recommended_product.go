package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// RecommendedProduct 推荐商品信息
type RecommendedProduct struct {

	// 推荐级别
	RecommendedLevel RecommendedProductRecommendedLevel `json:"recommended_level"`

	// 应用场景
	ApplicationScenarios RecommendedProductApplicationScenarios `json:"application_scenarios"`

	// 规格信息
	Flavors []RecommendFlavor `json:"flavors"`
}

func (o RecommendedProduct) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecommendedProduct struct{}"
	}

	return strings.Join([]string{"RecommendedProduct", string(data)}, " ")
}

type RecommendedProductRecommendedLevel struct {
	value string
}

type RecommendedProductRecommendedLevelEnum struct {
	E_0 RecommendedProductRecommendedLevel
	E_1 RecommendedProductRecommendedLevel
	E_2 RecommendedProductRecommendedLevel
}

func GetRecommendedProductRecommendedLevelEnum() RecommendedProductRecommendedLevelEnum {
	return RecommendedProductRecommendedLevelEnum{
		E_0: RecommendedProductRecommendedLevel{
			value: "0",
		},
		E_1: RecommendedProductRecommendedLevel{
			value: "1",
		},
		E_2: RecommendedProductRecommendedLevel{
			value: "2",
		},
	}
}

func (c RecommendedProductRecommendedLevel) Value() string {
	return c.value
}

func (c RecommendedProductRecommendedLevel) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RecommendedProductRecommendedLevel) UnmarshalJSON(b []byte) error {
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

type RecommendedProductApplicationScenarios struct {
	value string
}

type RecommendedProductApplicationScenariosEnum struct {
	E_0 RecommendedProductApplicationScenarios
	E_1 RecommendedProductApplicationScenarios
	E_2 RecommendedProductApplicationScenarios
}

func GetRecommendedProductApplicationScenariosEnum() RecommendedProductApplicationScenariosEnum {
	return RecommendedProductApplicationScenariosEnum{
		E_0: RecommendedProductApplicationScenarios{
			value: "0",
		},
		E_1: RecommendedProductApplicationScenarios{
			value: "1",
		},
		E_2: RecommendedProductApplicationScenarios{
			value: "2",
		},
	}
}

func (c RecommendedProductApplicationScenarios) Value() string {
	return c.value
}

func (c RecommendedProductApplicationScenarios) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RecommendedProductApplicationScenarios) UnmarshalJSON(b []byte) error {
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
