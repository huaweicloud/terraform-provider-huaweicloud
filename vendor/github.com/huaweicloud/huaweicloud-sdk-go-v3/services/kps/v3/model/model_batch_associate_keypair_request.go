package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchAssociateKeypairRequest Request Object
type BatchAssociateKeypairRequest struct {
	Body *BatchAssociateKeypairRequestBody `json:"body,omitempty"`
}

func (o BatchAssociateKeypairRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchAssociateKeypairRequest struct{}"
	}

	return strings.Join([]string{"BatchAssociateKeypairRequest", string(data)}, " ")
}
