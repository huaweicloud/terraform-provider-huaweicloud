package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateQualityEnhanceTemplateRequest struct {
	Body *QualityEnhanceTemplate `json:"body,omitempty"`
}

func (o CreateQualityEnhanceTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateQualityEnhanceTemplateRequest struct{}"
	}

	return strings.Join([]string{"CreateQualityEnhanceTemplateRequest", string(data)}, " ")
}
