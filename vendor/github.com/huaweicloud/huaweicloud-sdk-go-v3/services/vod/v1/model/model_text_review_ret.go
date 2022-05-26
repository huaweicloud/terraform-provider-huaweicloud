package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 文本检测结果
type TextReviewRet struct {

	// 检测结果是否通过。  取值如下： - block：包含敏感信息，不通过。 - pass：不包含敏感信息，通过。 - review：需要人工复检。
	Suggestion TextReviewRetSuggestion `json:"suggestion"`

	// 涉政敏感词列表
	Politics *string `json:"politics,omitempty"`

	// 涉黄敏感词列表
	Porn *string `json:"porn,omitempty"`

	// 辱骂敏感词列表
	Abuse *string `json:"abuse,omitempty"`
}

func (o TextReviewRet) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TextReviewRet struct{}"
	}

	return strings.Join([]string{"TextReviewRet", string(data)}, " ")
}

type TextReviewRetSuggestion struct {
	value string
}

type TextReviewRetSuggestionEnum struct {
	BLOCK  TextReviewRetSuggestion
	PASS   TextReviewRetSuggestion
	REVIEW TextReviewRetSuggestion
}

func GetTextReviewRetSuggestionEnum() TextReviewRetSuggestionEnum {
	return TextReviewRetSuggestionEnum{
		BLOCK: TextReviewRetSuggestion{
			value: "block",
		},
		PASS: TextReviewRetSuggestion{
			value: "pass",
		},
		REVIEW: TextReviewRetSuggestion{
			value: "review",
		},
	}
}

func (c TextReviewRetSuggestion) Value() string {
	return c.value
}

func (c TextReviewRetSuggestion) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TextReviewRetSuggestion) UnmarshalJSON(b []byte) error {
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
