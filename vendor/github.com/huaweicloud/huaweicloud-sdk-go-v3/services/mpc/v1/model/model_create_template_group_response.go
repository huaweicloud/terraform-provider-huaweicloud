package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateTemplateGroupResponse struct {
	TemplateGroup  *TemplateGroup `json:"template_group,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o CreateTemplateGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTemplateGroupResponse struct{}"
	}

	return strings.Join([]string{"CreateTemplateGroupResponse", string(data)}, " ")
}
