package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListQualityEnhanceDefaultTemplateRequest struct {
}

func (o ListQualityEnhanceDefaultTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListQualityEnhanceDefaultTemplateRequest struct{}"
	}

	return strings.Join([]string{"ListQualityEnhanceDefaultTemplateRequest", string(data)}, " ")
}
