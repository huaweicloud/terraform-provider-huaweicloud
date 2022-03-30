package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 计费类型信息，支持包年包月和按需，默认为按需。
type MysqlChargeInfo struct {
	// 计费模式。  取值范围：  - prePaid：预付费，即包年/包月。 - postPaid：后付费，即按需付费。

	ChargeMode MysqlChargeInfoChargeMode `json:"charge_mode"`
	// 订购周期类型。  取值范围：  - month：包月。 - year：包年。  说明：“charge_mode”为“prePaid”时生效，且为必选值。

	PeriodType *MysqlChargeInfoPeriodType `json:"period_type,omitempty"`
	// “charge_mode”为“prePaid”时生效，且为必选值，指定订购的时间。  取值范围：  当“period_type”为“month”时，取值为1~9。 当“period_type”为“year”时，取值为1~3。

	PeriodNum *int32 `json:"period_num,omitempty"`
	// 创建包周期实例时可指定，表示是否自动续订，续订的周期和原周期相同，且续订时会自动支付。  - true，为自动续订。 - false，为不自动续订，默认该方式。

	IsAutoRenew *string `json:"is_auto_renew,omitempty"`
	// 创建包周期时可指定，表示是否自动从客户的账户中支付，此字段不影响自动续订的支付方式。  - true，为自动支付，默认该方式。 - false，为手动支付。

	IsAutoPay *string `json:"is_auto_pay,omitempty"`
}

func (o MysqlChargeInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlChargeInfo struct{}"
	}

	return strings.Join([]string{"MysqlChargeInfo", string(data)}, " ")
}

type MysqlChargeInfoChargeMode struct {
	value string
}

type MysqlChargeInfoChargeModeEnum struct {
	PRE_PAID  MysqlChargeInfoChargeMode
	POST_PAID MysqlChargeInfoChargeMode
}

func GetMysqlChargeInfoChargeModeEnum() MysqlChargeInfoChargeModeEnum {
	return MysqlChargeInfoChargeModeEnum{
		PRE_PAID: MysqlChargeInfoChargeMode{
			value: "prePaid",
		},
		POST_PAID: MysqlChargeInfoChargeMode{
			value: "postPaid",
		},
	}
}

func (c MysqlChargeInfoChargeMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MysqlChargeInfoChargeMode) UnmarshalJSON(b []byte) error {
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

type MysqlChargeInfoPeriodType struct {
	value string
}

type MysqlChargeInfoPeriodTypeEnum struct {
	MONTH MysqlChargeInfoPeriodType
	YEAR  MysqlChargeInfoPeriodType
}

func GetMysqlChargeInfoPeriodTypeEnum() MysqlChargeInfoPeriodTypeEnum {
	return MysqlChargeInfoPeriodTypeEnum{
		MONTH: MysqlChargeInfoPeriodType{
			value: "month",
		},
		YEAR: MysqlChargeInfoPeriodType{
			value: "year",
		},
	}
}

func (c MysqlChargeInfoPeriodType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MysqlChargeInfoPeriodType) UnmarshalJSON(b []byte) error {
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
