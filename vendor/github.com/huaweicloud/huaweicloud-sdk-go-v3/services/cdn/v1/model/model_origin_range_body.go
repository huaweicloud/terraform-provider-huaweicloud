package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type OriginRangeBody struct {

	// range状态（\"off\"/\"on\"）
	RangeStatus *OriginRangeBodyRangeStatus `json:"range_status,omitempty"`

	// 加速域名id。
	DomainId *string `json:"domain_id,omitempty"`
}

func (o OriginRangeBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OriginRangeBody struct{}"
	}

	return strings.Join([]string{"OriginRangeBody", string(data)}, " ")
}

type OriginRangeBodyRangeStatus struct {
	value string
}

type OriginRangeBodyRangeStatusEnum struct {
	OFF OriginRangeBodyRangeStatus
	ON  OriginRangeBodyRangeStatus
}

func GetOriginRangeBodyRangeStatusEnum() OriginRangeBodyRangeStatusEnum {
	return OriginRangeBodyRangeStatusEnum{
		OFF: OriginRangeBodyRangeStatus{
			value: "off",
		},
		ON: OriginRangeBodyRangeStatus{
			value: "on",
		},
	}
}

func (c OriginRangeBodyRangeStatus) Value() string {
	return c.value
}

func (c OriginRangeBodyRangeStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *OriginRangeBodyRangeStatus) UnmarshalJSON(b []byte) error {
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
