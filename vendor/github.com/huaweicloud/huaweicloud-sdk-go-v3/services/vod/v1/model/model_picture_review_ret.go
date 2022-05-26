package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 图片审核结果
type PictureReviewRet struct {

	// 检测结果是否通过。  取值如下： - block：包含敏感信息，不通过。 - pass：不包含敏感信息，通过。 - review：需要人工复检。
	Suggestion *PictureReviewRetSuggestion `json:"suggestion,omitempty"`

	// 截图在视频中的时间偏移值。封面不涉及此字段  单位：秒。
	Offset *int32 `json:"offset,omitempty"`

	// 对应截图/封面的访问URL。
	Url string `json:"url"`

	// 政治因素审核结果。
	Politics *[]ReviewDetail `json:"politics,omitempty"`

	// 暴恐元素审核结果。
	Terrorism *[]ReviewDetail `json:"terrorism,omitempty"`

	// 涉黄内容审核结果。
	Porn *[]ReviewDetail `json:"porn,omitempty"`
}

func (o PictureReviewRet) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PictureReviewRet struct{}"
	}

	return strings.Join([]string{"PictureReviewRet", string(data)}, " ")
}

type PictureReviewRetSuggestion struct {
	value string
}

type PictureReviewRetSuggestionEnum struct {
	BLOCK  PictureReviewRetSuggestion
	PASS   PictureReviewRetSuggestion
	REVIEW PictureReviewRetSuggestion
}

func GetPictureReviewRetSuggestionEnum() PictureReviewRetSuggestionEnum {
	return PictureReviewRetSuggestionEnum{
		BLOCK: PictureReviewRetSuggestion{
			value: "block",
		},
		PASS: PictureReviewRetSuggestion{
			value: "pass",
		},
		REVIEW: PictureReviewRetSuggestion{
			value: "review",
		},
	}
}

func (c PictureReviewRetSuggestion) Value() string {
	return c.value
}

func (c PictureReviewRetSuggestion) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PictureReviewRetSuggestion) UnmarshalJSON(b []byte) error {
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
