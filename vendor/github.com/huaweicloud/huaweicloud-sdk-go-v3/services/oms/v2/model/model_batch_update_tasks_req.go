package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// BatchUpdateTasksReq This is a auto create Body Object
type BatchUpdateTasksReq struct {

	// 迁移任务组ID，表示批量更新该任务组下所有任务。 group_id和ids为二选一参数，不可同时存在或同时缺失。
	GroupId *string `json:"group_id,omitempty"`

	// 迁移任务id数组，包含所有需要批量更新操作的任务id。 group_id和ids为二选一参数，不可同时存在或同时缺失。
	Ids *[]int64 `json:"ids,omitempty"`

	// 配置流量控制策略。数组中一个元素对应一个时段的最大带宽，最多允许5个时段，且时段不能重叠。
	BandwidthPolicy []BandwidthPolicyDto `json:"bandwidth_policy"`

	// 任务优先级配置，存在高中低三个优先级档次，限制仅在等待中、已暂停、已失败的任务进行修改 HIGH：高优先级 MEDIUM：中优先级 LOW：低优先级
	TaskPriority *BatchUpdateTasksReqTaskPriority `json:"task_priority,omitempty"`
}

func (o BatchUpdateTasksReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchUpdateTasksReq struct{}"
	}

	return strings.Join([]string{"BatchUpdateTasksReq", string(data)}, " ")
}

type BatchUpdateTasksReqTaskPriority struct {
	value string
}

type BatchUpdateTasksReqTaskPriorityEnum struct {
	HIGH   BatchUpdateTasksReqTaskPriority
	MEDIUM BatchUpdateTasksReqTaskPriority
	LOW    BatchUpdateTasksReqTaskPriority
}

func GetBatchUpdateTasksReqTaskPriorityEnum() BatchUpdateTasksReqTaskPriorityEnum {
	return BatchUpdateTasksReqTaskPriorityEnum{
		HIGH: BatchUpdateTasksReqTaskPriority{
			value: "HIGH",
		},
		MEDIUM: BatchUpdateTasksReqTaskPriority{
			value: "MEDIUM",
		},
		LOW: BatchUpdateTasksReqTaskPriority{
			value: "LOW",
		},
	}
}

func (c BatchUpdateTasksReqTaskPriority) Value() string {
	return c.value
}

func (c BatchUpdateTasksReqTaskPriority) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BatchUpdateTasksReqTaskPriority) UnmarshalJSON(b []byte) error {
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
