package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type RunRecordRequest struct {

	// 操作行为。 取值如下： - START：对指定流开始录制，必须在直播流已经推送情况下才能正常启动，使用此命令启动录制的直播流如果发生了断流且超出断流时长，就会停止录制，并且重新推流后不会自动启动录制。 - STOP：对指定流停止录制。
	Action RunRecordRequestAction `json:"action"`

	Body *RecordControlInfo `json:"body,omitempty"`
}

func (o RunRecordRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RunRecordRequest struct{}"
	}

	return strings.Join([]string{"RunRecordRequest", string(data)}, " ")
}

type RunRecordRequestAction struct {
	value string
}

type RunRecordRequestActionEnum struct {
	START RunRecordRequestAction
	STOP  RunRecordRequestAction
}

func GetRunRecordRequestActionEnum() RunRecordRequestActionEnum {
	return RunRecordRequestActionEnum{
		START: RunRecordRequestAction{
			value: "START",
		},
		STOP: RunRecordRequestAction{
			value: "STOP",
		},
	}
}

func (c RunRecordRequestAction) Value() string {
	return c.value
}

func (c RunRecordRequestAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RunRecordRequestAction) UnmarshalJSON(b []byte) error {
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
