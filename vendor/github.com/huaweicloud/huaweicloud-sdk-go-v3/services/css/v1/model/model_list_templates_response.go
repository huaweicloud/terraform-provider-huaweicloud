package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTemplatesResponse Response Object
type ListTemplatesResponse struct {

	// 系统模板列表。
	SystemTemplates *[]SystemTemplates `json:"systemTemplates,omitempty"`

	// 自定义模板列表。
	CustomTemplates *[]CustomTemplates `json:"customTemplates,omitempty"`
	HttpStatusCode  int                `json:"-"`
}

func (o ListTemplatesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTemplatesResponse struct{}"
	}

	return strings.Join([]string{"ListTemplatesResponse", string(data)}, " ")
}
