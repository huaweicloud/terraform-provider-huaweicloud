package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ShowSinkTaskDetailRequest struct {

	// 实例转储ID。 请参考[实例生命周期][查询实例]接口返回的数据。
	ConnectorId string `json:"connector_id"`

	// 转储任务ID。
	TaskId string `json:"task_id"`

	// 是否包含topic信息。默认是false。
	TopicInfo *ShowSinkTaskDetailRequestTopicInfo `json:"topic-info,omitempty"`
}

func (o ShowSinkTaskDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowSinkTaskDetailRequest struct{}"
	}

	return strings.Join([]string{"ShowSinkTaskDetailRequest", string(data)}, " ")
}

type ShowSinkTaskDetailRequestTopicInfo struct {
	value string
}

type ShowSinkTaskDetailRequestTopicInfoEnum struct {
	TRUE  ShowSinkTaskDetailRequestTopicInfo
	FALSE ShowSinkTaskDetailRequestTopicInfo
}

func GetShowSinkTaskDetailRequestTopicInfoEnum() ShowSinkTaskDetailRequestTopicInfoEnum {
	return ShowSinkTaskDetailRequestTopicInfoEnum{
		TRUE: ShowSinkTaskDetailRequestTopicInfo{
			value: "true",
		},
		FALSE: ShowSinkTaskDetailRequestTopicInfo{
			value: "false",
		},
	}
}

func (c ShowSinkTaskDetailRequestTopicInfo) Value() string {
	return c.value
}

func (c ShowSinkTaskDetailRequestTopicInfo) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowSinkTaskDetailRequestTopicInfo) UnmarshalJSON(b []byte) error {
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
