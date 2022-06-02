package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 审核信息数组。  > 仅当审核成功后才能查询到此信息，未审核、正在审核以及审核失败时，无此字段信息。
type ReviewInfo struct {

	// 检测结果是否通过。  取值如下： - block：包含敏感信息，不通过。 - pass：不包含敏感信息，通过。 - review：需要人工复检。  > 当同时检测多个场景时，suggestion的值以最可能包含敏感信息的场景为准。即任一场景出现了block则总的suggestion为block，所有场景都pass时suggestion为pass，这两种情况之外则一定有场景需要review，此时suggestion为review。
	Suggestion *ReviewInfoSuggestion `json:"suggestion,omitempty"`

	Text *TextReviewRet `json:"text,omitempty"`

	// 封面检测结果。
	Cover *[]PictureReviewRet `json:"cover,omitempty"`

	// 视频检测结果。
	Video *[]PictureReviewRet `json:"video,omitempty"`

	// 执行情况描述。
	ExecDesc *string `json:"exec_desc,omitempty"`

	// 审核状态。  取值如下： - UN_REVIEW：未审核 - REVIEWING：审核中 - REVIEW_SUSPICIOUS：审核可疑，需要人工审核 - REVIEW_PASSED：审核通过 - REVIEW_FAILED：审核失败。 - REVIEW_BLOCKED：已屏蔽。
	ReviewStatus string `json:"review_status"`
}

func (o ReviewInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReviewInfo struct{}"
	}

	return strings.Join([]string{"ReviewInfo", string(data)}, " ")
}

type ReviewInfoSuggestion struct {
	value string
}

type ReviewInfoSuggestionEnum struct {
	BLOCK  ReviewInfoSuggestion
	PASS   ReviewInfoSuggestion
	REVIEW ReviewInfoSuggestion
}

func GetReviewInfoSuggestionEnum() ReviewInfoSuggestionEnum {
	return ReviewInfoSuggestionEnum{
		BLOCK: ReviewInfoSuggestion{
			value: "block",
		},
		PASS: ReviewInfoSuggestion{
			value: "pass",
		},
		REVIEW: ReviewInfoSuggestion{
			value: "review",
		},
	}
}

func (c ReviewInfoSuggestion) Value() string {
	return c.value
}

func (c ReviewInfoSuggestion) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ReviewInfoSuggestion) UnmarshalJSON(b []byte) error {
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
