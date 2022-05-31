package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListWatermarkTemplateRequest struct {

	// 水印模板配置id，一次最多10个。
	Id *[]string `json:"id,omitempty"`

	// 分页编号。  默认为0。指定id时该参数无效。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。  默认为10，范围[1,100]。指定id时该参数无效。
	Size *int32 `json:"size,omitempty"`
}

func (o ListWatermarkTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListWatermarkTemplateRequest struct{}"
	}

	return strings.Join([]string{"ListWatermarkTemplateRequest", string(data)}, " ")
}
