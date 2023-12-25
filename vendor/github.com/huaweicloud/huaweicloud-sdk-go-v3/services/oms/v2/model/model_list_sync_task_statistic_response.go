package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListSyncTaskStatisticResponse Response Object
type ListSyncTaskStatisticResponse struct {

	// 同步任务id
	SyncTaskId *string `json:"sync_task_id,omitempty"`

	// 统计结果时间间隔说明描述 FIVE_MINUTES：5分钟 ONE_HOUR：1小时
	StatisticTimeType *ListSyncTaskStatisticResponseStatisticTimeType `json:"statistic_time_type,omitempty"`

	// 查询的同步任务统计结果集
	StatisticDatas *[]StatisticTypeData `json:"statistic_datas,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o ListSyncTaskStatisticResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSyncTaskStatisticResponse struct{}"
	}

	return strings.Join([]string{"ListSyncTaskStatisticResponse", string(data)}, " ")
}

type ListSyncTaskStatisticResponseStatisticTimeType struct {
	value string
}

type ListSyncTaskStatisticResponseStatisticTimeTypeEnum struct {
	FIVE_MINUTES ListSyncTaskStatisticResponseStatisticTimeType
	ONE_HOUR     ListSyncTaskStatisticResponseStatisticTimeType
}

func GetListSyncTaskStatisticResponseStatisticTimeTypeEnum() ListSyncTaskStatisticResponseStatisticTimeTypeEnum {
	return ListSyncTaskStatisticResponseStatisticTimeTypeEnum{
		FIVE_MINUTES: ListSyncTaskStatisticResponseStatisticTimeType{
			value: "FIVE_MINUTES",
		},
		ONE_HOUR: ListSyncTaskStatisticResponseStatisticTimeType{
			value: "ONE_HOUR",
		},
	}
}

func (c ListSyncTaskStatisticResponseStatisticTimeType) Value() string {
	return c.value
}

func (c ListSyncTaskStatisticResponseStatisticTimeType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListSyncTaskStatisticResponseStatisticTimeType) UnmarshalJSON(b []byte) error {
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
