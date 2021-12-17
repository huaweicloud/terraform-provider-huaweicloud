package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// This is a auto create Body Object
type CreateTemplateReq struct {
	Template *TemplateRequest `json:"template"`
}

func (o CreateTemplateReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTemplateReq struct{}"
	}

	return strings.Join([]string{"CreateTemplateReq", string(data)}, " ")
}
