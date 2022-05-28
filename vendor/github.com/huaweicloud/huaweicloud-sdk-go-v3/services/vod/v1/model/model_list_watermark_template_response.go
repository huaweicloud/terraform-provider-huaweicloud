package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListWatermarkTemplateResponse struct {

	// 水印模板信息。
	Templates *[]WatermarkTemplate `json:"templates,omitempty"`

	// 总数。
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListWatermarkTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListWatermarkTemplateResponse struct{}"
	}

	return strings.Join([]string{"ListWatermarkTemplateResponse", string(data)}, " ")
}
