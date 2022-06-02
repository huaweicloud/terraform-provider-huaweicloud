package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ShowMetricsDataRequest struct {

	// 用于对查询到的监控数据进行断点插值，默认值为-1。 -1：断点处使用-1进行表示。 0 ：断点处使用0进行表示。 null：断点处使用null进行表示。 average：断点处使用前后邻近的有效数据的平均值进行表示，如果不存在有效数据则使用null进行表示。
	FillValue *ShowMetricsDataRequestFillValue `json:"fillValue,omitempty"`

	Body *QueryMetricDataParam `json:"body,omitempty"`
}

func (o ShowMetricsDataRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMetricsDataRequest struct{}"
	}

	return strings.Join([]string{"ShowMetricsDataRequest", string(data)}, " ")
}

type ShowMetricsDataRequestFillValue struct {
	value string
}

type ShowMetricsDataRequestFillValueEnum struct {
	E_1     ShowMetricsDataRequestFillValue
	E_0     ShowMetricsDataRequestFillValue
	NULL    ShowMetricsDataRequestFillValue
	AVERAGE ShowMetricsDataRequestFillValue
}

func GetShowMetricsDataRequestFillValueEnum() ShowMetricsDataRequestFillValueEnum {
	return ShowMetricsDataRequestFillValueEnum{
		E_1: ShowMetricsDataRequestFillValue{
			value: "-1",
		},
		E_0: ShowMetricsDataRequestFillValue{
			value: "0",
		},
		NULL: ShowMetricsDataRequestFillValue{
			value: "null",
		},
		AVERAGE: ShowMetricsDataRequestFillValue{
			value: "average",
		},
	}
}

func (c ShowMetricsDataRequestFillValue) Value() string {
	return c.value
}

func (c ShowMetricsDataRequestFillValue) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowMetricsDataRequestFillValue) UnmarshalJSON(b []byte) error {
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
