package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type QueryTranscodingsTaskResponse struct {

	// 任务ID。
	TaskId *string `json:"task_id,omitempty"`

	// 任务执行状态。
	Status *QueryTranscodingsTaskResponseStatus `json:"status,omitempty"`

	// 转码任务启动时间
	CreateTime *string `json:"create_time,omitempty"`

	// 转码任务结束时间
	EndTime *string `json:"end_time,omitempty"`

	// 转码任务对应的转码模板ID
	TransTemplateId *[]int32 `json:"trans_template_id,omitempty"`

	Input *ObsObjInfo `json:"input,omitempty"`

	Output *ObsObjInfo `json:"output,omitempty"`

	// 转码生成的文件名，数组类型，可能包含多个，包含截图文件名。
	OutputFileName *[]string `json:"output_file_name,omitempty"`

	// 用户自定义数据。
	UserData *string `json:"user_data,omitempty"`

	// 转码任务的返回码。
	ErrorCode *string `json:"error_code,omitempty"`

	// 转码任务描述，当转码出现异常时，此字段为异常的原因。
	Description *string `json:"description,omitempty"`

	// 转码成功，但音频采样率过低时的提示。
	Tips *string `json:"tips,omitempty"`

	TranscodeDetail *TranscodeDetail `json:"transcode_detail,omitempty"`

	ThumbnailOutput *ObsObjInfo `json:"thumbnail_output,omitempty"`

	// 截图压缩包名。
	ThumbnailOutputname *string `json:"thumbnail_outputname,omitempty"`

	// 截图文件信息。
	PicInfo *[]PicInfo `json:"pic_info,omitempty"`

	// 转码参数。  若同时设置“trans_template_id”和此参数，则优先使用此参数进行转码。
	AvParameters *[]AvParameters `json:"av_parameters,omitempty"`
}

func (o QueryTranscodingsTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QueryTranscodingsTaskResponse struct{}"
	}

	return strings.Join([]string{"QueryTranscodingsTaskResponse", string(data)}, " ")
}

type QueryTranscodingsTaskResponseStatus struct {
	value string
}

type QueryTranscodingsTaskResponseStatusEnum struct {
	NO_TASK          QueryTranscodingsTaskResponseStatus
	WAITING          QueryTranscodingsTaskResponseStatus
	TRANSCODING      QueryTranscodingsTaskResponseStatus
	SUCCEEDED        QueryTranscodingsTaskResponseStatus
	FAILED           QueryTranscodingsTaskResponseStatus
	CANCELED         QueryTranscodingsTaskResponseStatus
	NEED_TO_BE_AUDIT QueryTranscodingsTaskResponseStatus
}

func GetQueryTranscodingsTaskResponseStatusEnum() QueryTranscodingsTaskResponseStatusEnum {
	return QueryTranscodingsTaskResponseStatusEnum{
		NO_TASK: QueryTranscodingsTaskResponseStatus{
			value: "NO_TASK",
		},
		WAITING: QueryTranscodingsTaskResponseStatus{
			value: "WAITING",
		},
		TRANSCODING: QueryTranscodingsTaskResponseStatus{
			value: "TRANSCODING",
		},
		SUCCEEDED: QueryTranscodingsTaskResponseStatus{
			value: "SUCCEEDED",
		},
		FAILED: QueryTranscodingsTaskResponseStatus{
			value: "FAILED",
		},
		CANCELED: QueryTranscodingsTaskResponseStatus{
			value: "CANCELED",
		},
		NEED_TO_BE_AUDIT: QueryTranscodingsTaskResponseStatus{
			value: "NEED_TO_BE_AUDIT",
		},
	}
}

func (c QueryTranscodingsTaskResponseStatus) Value() string {
	return c.value
}

func (c QueryTranscodingsTaskResponseStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *QueryTranscodingsTaskResponseStatus) UnmarshalJSON(b []byte) error {
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
