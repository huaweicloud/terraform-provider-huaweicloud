package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MultiTaskInfo struct {

	// 转码模板ID。
	TemplateId *int32 `json:"template_id,omitempty"`

	Error *ErrorResponse `json:"error,omitempty"`

	OutputFile *SourceInfo `json:"output_file,omitempty"`
}

func (o MultiTaskInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MultiTaskInfo struct{}"
	}

	return strings.Join([]string{"MultiTaskInfo", string(data)}, " ")
}
