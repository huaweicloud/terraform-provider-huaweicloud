package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListSyncTaskStatisticRequest Request Object
type ListSyncTaskStatisticRequest struct {

	// 同步任务ID。
	SyncTaskId string `json:"sync_task_id"`

	// 统计数据类型： 多类型查询用‘,’分割； REQUEST：接收同步请求对象数 SUCCESS：同步成功对象数 FAILURE：同步失败对象数 SKIP：同步跳过对象数 SIZE：同步成功对象容量(Byte)
	DataType ListSyncTaskStatisticRequestDataType `json:"data_type"`

	// 查询开始时间
	StartTime string `json:"start_time"`

	// 查询开始时间
	EndTime string `json:"end_time"`
}

func (o ListSyncTaskStatisticRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSyncTaskStatisticRequest struct{}"
	}

	return strings.Join([]string{"ListSyncTaskStatisticRequest", string(data)}, " ")
}

type ListSyncTaskStatisticRequestDataType struct {
	value string
}

type ListSyncTaskStatisticRequestDataTypeEnum struct {
	REQUEST ListSyncTaskStatisticRequestDataType
	SUCCESS ListSyncTaskStatisticRequestDataType
	FAILURE ListSyncTaskStatisticRequestDataType
	SKIP    ListSyncTaskStatisticRequestDataType
	SIZE    ListSyncTaskStatisticRequestDataType
}

func GetListSyncTaskStatisticRequestDataTypeEnum() ListSyncTaskStatisticRequestDataTypeEnum {
	return ListSyncTaskStatisticRequestDataTypeEnum{
		REQUEST: ListSyncTaskStatisticRequestDataType{
			value: "REQUEST",
		},
		SUCCESS: ListSyncTaskStatisticRequestDataType{
			value: "SUCCESS",
		},
		FAILURE: ListSyncTaskStatisticRequestDataType{
			value: "FAILURE",
		},
		SKIP: ListSyncTaskStatisticRequestDataType{
			value: "SKIP",
		},
		SIZE: ListSyncTaskStatisticRequestDataType{
			value: "SIZE",
		},
	}
}

func (c ListSyncTaskStatisticRequestDataType) Value() string {
	return c.value
}

func (c ListSyncTaskStatisticRequestDataType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListSyncTaskStatisticRequestDataType) UnmarshalJSON(b []byte) error {
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
