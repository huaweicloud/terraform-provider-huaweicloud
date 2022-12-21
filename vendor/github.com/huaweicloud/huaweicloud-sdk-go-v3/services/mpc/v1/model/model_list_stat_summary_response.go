package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Response Object
type ListStatSummaryResponse struct {

	// 统计概览信息
	Summary *[]StatSummary `json:"summary,omitempty"`

	// 该指标的总值，精确到小数点后两位。
	Total *float32 `json:"total,omitempty"`

	// 统计类型。取值如下： - video_duration, 转码片源时长统计，单位：分钟。 - remux_file_duration，转封装片源时长统计，单位：分钟。 - transcode_task_number，转码次数统计，单位：次。 - transcode_duration，转码耗时时长统计，单位：分钟。
	StatType       *ListStatSummaryResponseStatType `json:"stat_type,omitempty"`
	HttpStatusCode int                              `json:"-"`
}

func (o ListStatSummaryResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListStatSummaryResponse struct{}"
	}

	return strings.Join([]string{"ListStatSummaryResponse", string(data)}, " ")
}

type ListStatSummaryResponseStatType struct {
	value string
}

type ListStatSummaryResponseStatTypeEnum struct {
	VIDEO_DURATION        ListStatSummaryResponseStatType
	REMUX_FILE_DURATION   ListStatSummaryResponseStatType
	TRANSCODE_TASK_NUMBER ListStatSummaryResponseStatType
	TRANSCODE_DURATION    ListStatSummaryResponseStatType
}

func GetListStatSummaryResponseStatTypeEnum() ListStatSummaryResponseStatTypeEnum {
	return ListStatSummaryResponseStatTypeEnum{
		VIDEO_DURATION: ListStatSummaryResponseStatType{
			value: "video_duration",
		},
		REMUX_FILE_DURATION: ListStatSummaryResponseStatType{
			value: "remux_file_duration",
		},
		TRANSCODE_TASK_NUMBER: ListStatSummaryResponseStatType{
			value: "transcode_task_number",
		},
		TRANSCODE_DURATION: ListStatSummaryResponseStatType{
			value: "transcode_duration",
		},
	}
}

func (c ListStatSummaryResponseStatType) Value() string {
	return c.value
}

func (c ListStatSummaryResponseStatType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListStatSummaryResponseStatType) UnmarshalJSON(b []byte) error {
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
