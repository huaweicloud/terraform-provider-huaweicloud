package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateTransTemplateResponse struct {

	// 自定义转码模板编号。
	TemplateId     *int32 `json:"template_id,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o CreateTransTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTransTemplateResponse struct{}"
	}

	return strings.Join([]string{"CreateTransTemplateResponse", string(data)}, " ")
}
