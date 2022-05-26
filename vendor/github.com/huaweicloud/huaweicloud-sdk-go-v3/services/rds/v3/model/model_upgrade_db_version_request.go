package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type UpgradeDbVersionRequest struct {

	// 语言
	XLanguage *UpgradeDbVersionRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *CustomerUpgradeDatabaseVersionReq `json:"body,omitempty"`
}

func (o UpgradeDbVersionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeDbVersionRequest struct{}"
	}

	return strings.Join([]string{"UpgradeDbVersionRequest", string(data)}, " ")
}

type UpgradeDbVersionRequestXLanguage struct {
	value string
}

type UpgradeDbVersionRequestXLanguageEnum struct {
	ZH_CN UpgradeDbVersionRequestXLanguage
	EN_US UpgradeDbVersionRequestXLanguage
}

func GetUpgradeDbVersionRequestXLanguageEnum() UpgradeDbVersionRequestXLanguageEnum {
	return UpgradeDbVersionRequestXLanguageEnum{
		ZH_CN: UpgradeDbVersionRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: UpgradeDbVersionRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c UpgradeDbVersionRequestXLanguage) Value() string {
	return c.value
}

func (c UpgradeDbVersionRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpgradeDbVersionRequestXLanguage) UnmarshalJSON(b []byte) error {
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
