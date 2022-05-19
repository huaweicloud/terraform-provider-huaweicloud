package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListTemplateResponse struct {

	// 转码模板
	TemplateArray *[]TemplateInfo `json:"template_array,omitempty"`

	// 转码模板总数
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTemplateResponse struct{}"
	}

	return strings.Join([]string{"ListTemplateResponse", string(data)}, " ")
}
