package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type MpeCallBackReq struct {

	// 任务类型
	TaskType *MpeCallBackReqTaskType `json:"task_type,omitempty"`

	// 任务ID。
	TaskId *string `json:"task_id,omitempty"`

	// 任务状态。
	Status *string `json:"status,omitempty"`

	// 任务完成进度百分比值。
	CompleteRatio *int32 `json:"complete_ratio,omitempty"`

	// 任务执行描述。
	Description *string `json:"description,omitempty"`

	MetaData *MpeMetaData `json:"meta_data,omitempty"`
}

func (o MpeCallBackReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MpeCallBackReq struct{}"
	}

	return strings.Join([]string{"MpeCallBackReq", string(data)}, " ")
}

type MpeCallBackReqTaskType struct {
	value string
}

type MpeCallBackReqTaskTypeEnum struct {
	CONCAT    MpeCallBackReqTaskType
	AUDIO     MpeCallBackReqTaskType
	CUT       MpeCallBackReqTaskType
	PARSE     MpeCallBackReqTaskType
	MD5       MpeCallBackReqTaskType
	SNAPSHOT  MpeCallBackReqTaskType
	REMUX     MpeCallBackReqTaskType
	ANIMATION MpeCallBackReqTaskType
}

func GetMpeCallBackReqTaskTypeEnum() MpeCallBackReqTaskTypeEnum {
	return MpeCallBackReqTaskTypeEnum{
		CONCAT: MpeCallBackReqTaskType{
			value: "CONCAT",
		},
		AUDIO: MpeCallBackReqTaskType{
			value: "AUDIO",
		},
		CUT: MpeCallBackReqTaskType{
			value: "CUT",
		},
		PARSE: MpeCallBackReqTaskType{
			value: "PARSE",
		},
		MD5: MpeCallBackReqTaskType{
			value: "MD5",
		},
		SNAPSHOT: MpeCallBackReqTaskType{
			value: "SNAPSHOT",
		},
		REMUX: MpeCallBackReqTaskType{
			value: "REMUX",
		},
		ANIMATION: MpeCallBackReqTaskType{
			value: "ANIMATION",
		},
	}
}

func (c MpeCallBackReqTaskType) Value() string {
	return c.value
}

func (c MpeCallBackReqTaskType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MpeCallBackReqTaskType) UnmarshalJSON(b []byte) error {
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
