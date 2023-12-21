package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// StatisticTypeData 查询同步任务统计结果
type StatisticTypeData struct {

	// 统计数据类型： REQUEST：请求对象数 SUCCESS：成功对象数 FAILURE：失败对象数 SKIP：跳过对象数 SIZE：对象容量(Byte)
	DataType *StatisticTypeDataDataType `json:"data_type,omitempty"`

	// 查询的同步任务统计结果集
	Data *[]StatisticData `json:"data,omitempty"`
}

func (o StatisticTypeData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StatisticTypeData struct{}"
	}

	return strings.Join([]string{"StatisticTypeData", string(data)}, " ")
}

type StatisticTypeDataDataType struct {
	value string
}

type StatisticTypeDataDataTypeEnum struct {
	REQUEST StatisticTypeDataDataType
	SUCCESS StatisticTypeDataDataType
	FAILURE StatisticTypeDataDataType
	SKIP    StatisticTypeDataDataType
	SIZE    StatisticTypeDataDataType
}

func GetStatisticTypeDataDataTypeEnum() StatisticTypeDataDataTypeEnum {
	return StatisticTypeDataDataTypeEnum{
		REQUEST: StatisticTypeDataDataType{
			value: "REQUEST",
		},
		SUCCESS: StatisticTypeDataDataType{
			value: "SUCCESS",
		},
		FAILURE: StatisticTypeDataDataType{
			value: "FAILURE",
		},
		SKIP: StatisticTypeDataDataType{
			value: "SKIP",
		},
		SIZE: StatisticTypeDataDataType{
			value: "SIZE",
		},
	}
}

func (c StatisticTypeDataDataType) Value() string {
	return c.value
}

func (c StatisticTypeDataDataType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *StatisticTypeDataDataType) UnmarshalJSON(b []byte) error {
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
