package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListQualityEnhanceDefaultTemplateResponse struct {

	// 任务列表
	TaskArray *[]QualityEnhanceTemplateInfo `json:"task_array,omitempty"`

	// 查询结果数量
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListQualityEnhanceDefaultTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListQualityEnhanceDefaultTemplateResponse struct{}"
	}

	return strings.Join([]string{"ListQualityEnhanceDefaultTemplateResponse", string(data)}, " ")
}
