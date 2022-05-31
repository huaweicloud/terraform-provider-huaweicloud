package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 审核结果
type ReviewDetail struct {

	// 置信度。  取值范围：[0,1]。
	Confidence string `json:"confidence"`

	// 每个检测结果的标签化说明。 - politics场景：label为对应的政治人物信息。 - terrorism场景： label为对应的暴恐元素（枪支、刀具、火灾等） 信息。 - porn场景：label为对应的涉黄元素（涉黄、性感等）信息。
	Label *string `json:"label,omitempty"`
}

func (o ReviewDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReviewDetail struct{}"
	}

	return strings.Join([]string{"ReviewDetail", string(data)}, " ")
}
