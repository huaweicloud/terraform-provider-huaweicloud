package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListTemplateRequest struct {

	// 自定义转码模板ID，最多10个
	TemplateId *[]int32 `json:"template_id,omitempty"`

	// 分页编号。查询指定“task_id”时，该参数无效。  默认值：0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。取值范围：[1,100]，指定template_id时该参数无效
	Size *int32 `json:"size,omitempty"`
}

func (o ListTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTemplateRequest struct{}"
	}

	return strings.Join([]string{"ListTemplateRequest", string(data)}, " ")
}
