package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListTemplatesResponse struct {
	// 模板个数

	Count *int32 `json:"count,omitempty"`
	// 模板信息

	Templates      *[]TemplateResponse `json:"templates,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ListTemplatesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTemplatesResponse struct{}"
	}

	return strings.Join([]string{"ListTemplatesResponse", string(data)}, " ")
}
