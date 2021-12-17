package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// This is a auto create Body Object
type UpdateTemplateReq struct {
	Template *TemplateRequest `json:"template,omitempty"`
}

func (o UpdateTemplateReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTemplateReq struct{}"
	}

	return strings.Join([]string{"UpdateTemplateReq", string(data)}, " ")
}
