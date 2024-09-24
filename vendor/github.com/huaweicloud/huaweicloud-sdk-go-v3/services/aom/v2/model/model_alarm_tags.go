package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AlarmTags struct {

	// 自动标签。
	AutoTags *[]string `json:"auto_tags,omitempty"`

	// 自定义标签。
	CustomTags *[]string `json:"custom_tags,omitempty"`

	// 告警标注。
	CustomAnnotations *[]string `json:"custom_annotations,omitempty"`
}

func (o AlarmTags) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AlarmTags struct{}"
	}

	return strings.Join([]string{"AlarmTags", string(data)}, " ")
}
