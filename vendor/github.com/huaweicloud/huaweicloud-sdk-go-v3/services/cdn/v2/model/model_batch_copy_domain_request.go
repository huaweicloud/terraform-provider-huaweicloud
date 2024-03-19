package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCopyDomainRequest Request Object
type BatchCopyDomainRequest struct {
	Body *BatchCopyDRequestBody `json:"body,omitempty"`
}

func (o BatchCopyDomainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCopyDomainRequest struct{}"
	}

	return strings.Join([]string{"BatchCopyDomainRequest", string(data)}, " ")
}
