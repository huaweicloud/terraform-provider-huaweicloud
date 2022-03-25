package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type AssociateKeypairRequest struct {
	Body *AssociateKeypairRequestBody `json:"body,omitempty"`
}

func (o AssociateKeypairRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociateKeypairRequest struct{}"
	}

	return strings.Join([]string{"AssociateKeypairRequest", string(data)}, " ")
}
