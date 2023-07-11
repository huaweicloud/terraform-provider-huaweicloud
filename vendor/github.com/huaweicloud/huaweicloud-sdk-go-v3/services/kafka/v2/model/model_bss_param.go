package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// BssParam 表示包周期计费模式的相关参数。  如果为空，则默认计费模式为按需计费；否则是包周期方式。
type BssParam struct {

	// 是否自动续订。  取值范围：   - true: 自动续订。   - false: 不自动续订。  默认不自动续订。
	IsAutoRenew *bool `json:"is_auto_renew,omitempty"`

	// 计费模式。  功能说明：付费方式。  取值范围：   - prePaid：预付费，即包年包月；   - postPaid：后付费，即按需付费；  默认为postPaid。
	ChargingMode *BssParamChargingMode `json:"charging_mode,omitempty"`

	// 下单订购后，是否自动从客户的账户中支付，而不需要客户手动去进行支付。  取值范围：   - true：是（自动支付）   - false：否（需要客户手动支付）  默认为手动支付。
	IsAutoPay *bool `json:"is_auto_pay,omitempty"`

	// 订购周期类型。  取值范围：   - month：月   - year：年  **chargingMode为prePaid时生效且为必选值。**
	PeriodType *BssParamPeriodType `json:"period_type,omitempty"`

	// 订购周期数。  取值范围：   - periodType=month（周期类型为月）时，取值为[1，9]；  - periodType=year（周期类型为年）时，取值为[1，3]；  **chargingMode为prePaid时生效且为必选值。**
	PeriodNum *int32 `json:"period_num,omitempty"`
}

func (o BssParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BssParam struct{}"
	}

	return strings.Join([]string{"BssParam", string(data)}, " ")
}

type BssParamChargingMode struct {
	value string
}

type BssParamChargingModeEnum struct {
	PRE_PAID  BssParamChargingMode
	POST_PAID BssParamChargingMode
}

func GetBssParamChargingModeEnum() BssParamChargingModeEnum {
	return BssParamChargingModeEnum{
		PRE_PAID: BssParamChargingMode{
			value: "prePaid",
		},
		POST_PAID: BssParamChargingMode{
			value: "postPaid",
		},
	}
}

func (c BssParamChargingMode) Value() string {
	return c.value
}

func (c BssParamChargingMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BssParamChargingMode) UnmarshalJSON(b []byte) error {
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

type BssParamPeriodType struct {
	value string
}

type BssParamPeriodTypeEnum struct {
	MONTH BssParamPeriodType
	YEAR  BssParamPeriodType
}

func GetBssParamPeriodTypeEnum() BssParamPeriodTypeEnum {
	return BssParamPeriodTypeEnum{
		MONTH: BssParamPeriodType{
			value: "month",
		},
		YEAR: BssParamPeriodType{
			value: "year",
		},
	}
}

func (c BssParamPeriodType) Value() string {
	return c.value
}

func (c BssParamPeriodType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BssParamPeriodType) UnmarshalJSON(b []byte) error {
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
